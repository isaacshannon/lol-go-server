function startUngank() {
    var streaming = false;
    var predicting = false;
    var reloading = false;

    var startbutton = null;
    var video = null;
    var overlay = null;
    var noOverlay = null;

    var userID = Math.floor(Math.random() * 10000000);
    var predictions = [[5, 5, "b"]];
    var heatmap = {};
    var sctx = null;
    var nsctx = null;

    const pause = time => new Promise(resolve => setTimeout(resolve, time))

    function handleButtonClick(ev) {
        const constraints = {
            audio: false,
            video: {
                facingMode: "environment"
            },
            width: 600,
        };

        if (reloading) {
            location.reload();
        }

        if (!streaming && !reloading) {
            if ('mediaDevices' in navigator && navigator.mediaDevices.getUserMedia) {
                navigator.mediaDevices.getUserMedia(constraints).then(handleStream).catch(function (err) {
                    logToServer(err);
                    startbutton.innerHTML = 'reset';
                    reloading = true;
                });
            }
            predicting = true;
            startbutton.innerHTML = 'reset';
            ev.preventDefault();
            return;
        }

        location.reload();
    }

    function handleStream(stream) {
        video.srcObject = stream;
        streaming = true;
        video.play()
        pause(1000).then(predictPositions);
    }

    function predictPositions() {
        logToServer({"name": "started prediction"});

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
                imgBase64: canvasData
            }
        }).done(function (d) {
            if (!("predictions" in d)){
                return;
            }
            if (d["predictions"].length === 0) {
                predictions = [0,0,"r"]
            } else {
                predictions = d["predictions"];
            }

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
        event["userID"] = userID
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
        '#FF000022',
        '#FF000033',
        '#FF000044',
        '#FF000055',
        '#FF000077'
    ];

    function drawSquares() {
        sctx.beginPath();
        var key;
        for (key in heatmap) {
            if (heatmap.hasOwnProperty(key)) {
                let x = key / 100;
                let y = key % 100;
                sctx.fillStyle = fills[heatmap[key]];
                sctx.fillRect(x * 30 - 15, y * 30 - 15, 30, 30);
            }
        }
        // alignment rectangle
        sctx.rect(20, 20, 260, 260);
        sctx.stroke();
    }

    overlay = document.getElementById('overlay');
    noOverlay = document.getElementById('noOverlay');
    sctx = overlay.getContext('2d');
    // alignment rectangle
    sctx.lineWidth = "1";
    sctx.strokeStyle = '#00FF00FF'
    sctx.rect(20, 20, 260, 260);
    sctx.stroke();
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

    function loadAndDrawImage(url, ctx){
        var image = new Image();
        image.onload = function()
        {
            ctx.drawImage(image, 0, 0);
        }
        image.src = url;
    }
    loadAndDrawImage("../static/map-holder.png", sctx);
}