function startUngank() {
    var width = 1;
    var height = 1;

    var streaming = false;
    var predicting = false;
    var reloading = false;

    var video = null;
    var startbutton = null;
    var overlay = null;
    var noOverlay = null;

    var userID = Math.floor(Math.random() * 10000000);
    var predictions = [[5, 5, "b"]];
    var heatmap = {};
    var sctx = null;
    var nsctx = null;

    const constraints = {
        audio: false,
        video: {
            facingMode: "environment"
        },
        width: 600,
    };

    const pause = time => new Promise(resolve => setTimeout(resolve, time))

    function handleButtonClick(ev) {
        logToServer({"name": "capture button clicked"});
        if (reloading) {
            logToServer({"name": "capture button error reload"});
            location.reload();
        }

        if (!streaming && !reloading) {
            logToServer({"name": "capture button starting stream"});
            if ('mediaDevices' in navigator && navigator.mediaDevices.getUserMedia) {
                navigator.mediaDevices.getUserMedia(constraints).then(handleStream).catch(function (err) {
                    logToServer(err);
                    startbutton.innerHTML = 'reset';
                    reloading = true;
                });
            }
            startbutton.innerHTML = 'start predictions';
            return;
        }

        if (streaming && !predicting) {
            logToServer({"name": "capture button starting predictions"});
            video.setAttribute('hidden', true);
            overlay.removeAttribute('hidden');
            predicting = true;
            startbutton.innerHTML = 'reset';
            video.play();
            video.setAttribute('width', width);
            video.setAttribute('height', height);
            pause(500).then(predictPositions);
            ev.preventDefault();
            return;
        }

        logToServer({"name": "capture button reload"});
        location.reload();
    }

    function handleStream(stream) {
        video.srcObject = stream;
        streaming = true;
        logToServer({"name": "stream started"})
    }

    function predictPositions() {
        logToServer({
            "name": "started prediction"
        });

        var canvasData = noOverlay.toDataURL('image/png');
        if (canvasData.length < 10000){
            // Canvas not initialized yet.
            pause(300).then(predictPositions);
            return;
        }
        $.ajax({
            type: "POST",
            url: "https://ungank.com/predict",
            data: {
                userID: userID,
                imgBase64: canvasData,
                x0: 1,
                x1: 1,
                y0: 1,
                y1: 1,
            }
        }).done(function (d) {
            if (!("predictions" in d) || d["predictions"].length === 0){
                return;
            }
            predictions = d["predictions"];
            var key;
            for (key in heatmap) {
                if (heatmap.hasOwnProperty(key) && heatmap[key] >= 1) {
                    heatmap[key] -= 1;
                }
            }
            for (let i = 0; i < predictions.length; i++) {
                let p = predictions[i];
                if (p[2] === "b"){
                    continue;
                }
                key = 100 * p[0] + p[1];
                if (p[3] > 0.2) {
                    heatmap[key] = 5;
                } else if (p[3] > 0.15 && heatmap[key] < 3) {
                    heatmap[key] = 3;
                } else if (p[3] > 0.12 && heatmap[key] < 2) {
                    heatmap[key] = 2;
                }
            }
            pause(300).then(predictPositions);
        });

    }

    function logToServer(event) {
        $.ajax({
            type: "POST",
            url: "https://ungank.com/log",
            data: event,
            async: true,
        })
    }


    console.log("starting up");
    video = document.getElementById('video');
    startbutton = document.getElementById('startbutton');
    startbutton.addEventListener('click', handleButtonClick, false);

    var fills = [
        '#00000000',
        '#99FF0011',
        '#FFFF0022',
        '#FF880033',
        '#FF330044',
        '#FF000077'
    ];

    function drawSquares() {
        sctx.beginPath();
        sctx.lineWidth = "1";
        var key;
        for (key in heatmap) {
            if (heatmap.hasOwnProperty(key)) {
                let x = key / 100;
                let y = key % 100;
                sctx.fillStyle = fills[heatmap[key]];
                sctx.fillRect(x * 30 - 15, y * 30 - 15, 30, 30);
            }
        }
    }

    overlay = document.getElementById('overlay');
    noOverlay = document.getElementById('noOverlay');
    sctx = overlay.getContext('2d');
    nsctx = noOverlay.getContext('2d');
    var i;
    video.addEventListener('play',
        function () {
            i = window.setInterval(
                function () {

                    var vWidth = video.videoWidth;
                    var vHeight = video.videoHeight;
                    if (vHeight > vWidth) {
                        var dy = (vHeight - vWidth) / 2;
                        sctx.drawImage(video, 0, dy, vWidth, vWidth, 0, 0, 300, 300);
                        nsctx.drawImage(video, 0, dy, vWidth, vWidth, 0, 0, 300, 300);
                    } else {
                        var dx = (vWidth - vHeight) / 2;
                        sctx.drawImage(video, dx, 0, vHeight, vHeight, 0, 0, 300, 300);
                        nsctx.drawImage(video, dx, 0, vHeight, vHeight, 0, 0, 300, 300);
                    }

                    drawSquares()

                }, 20);
        }, false);
}