function startUngank() {
    var width = 1;
    var height = 1;

    var streaming = false;
    var predicting = false;
    var reloading = false;

    var video = null;
    var startbutton = null;
    var overlay = null;
    var userID = Math.floor(Math.random() * 10000000);
    var predictions = [[5,5]];
    var sctx = null;

    const constraints = {
        audio: false,
        video: {
            facingMode: "environment"
        },
        width: 600,
    };

    function handleButtonClick(ev) {
        logToServer({"name":"capture button clicked"});
        if (reloading) {
            logToServer({"name":"capture button error reload"});
            location.reload();
        }

        if (!streaming && !reloading) {
            logToServer({"name":"capture button starting stream"});
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
            logToServer({"name":"capture button starting predictions"});
            video.setAttribute('hidden', true);
            overlay.removeAttribute('hidden');
            predicting = true;
            startbutton.innerHTML = 'reset';
            predictPositions();
            ev.preventDefault();
            return;
        }

        logToServer({"name":"capture button reload"});
        location.reload();
    }

    function handleStream(stream) {
        video.srcObject = stream;
        streaming = true;
        logToServer({"name": "stream started"})
    }

    function predictPositions() {
        height = video.videoHeight;
        width = video.videoWidth;

        video.setAttribute('width', width);
        video.setAttribute('height', height);

        logToServer({
            "name":"started prediction",
            "width": width,
            "height": height,
        });
        if (width && height) {
            var canvasData = overlay.toDataURL('image/png');
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
                logToServer({"prediction received": d["predictions"][0][0]});
                predictions = d["predictions"];
                predictPositions();

            });
        }
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

    function drawSquare(value, index, array) {
        sctx.beginPath();
        sctx.lineWidth = "1";
        sctx.fillStyle = '#FFFFFF55';
        sctx.fillRect(Number(value[0])*10, Number(value[1])*10, 10, 10);
        sctx.stroke();
    }

    overlay = document.getElementById('overlay');
    sctx = overlay.getContext('2d');
    var i;
    video.addEventListener('play',
        function() {
        i=window.setInterval(
            function() {

                var vWidth = video.width;
                var vHeight = video.height;
                if (vHeight > vWidth) {
                    var dy = (vHeight - vWidth) / 2;
                    sctx.drawImage(video, 0, dy, vWidth, vWidth, 0, 0, 300, 300)
                } else {
                    var dx = (vWidth - vHeight) / 2;
                    sctx.drawImage(video, dx, 0, vHeight, vHeight, 0, 0, 300, 300)
                }

                predictions.forEach(drawSquare);
            },20);
        },false);
}