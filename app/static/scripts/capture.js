function startUngank() {
    var width = 1;
    var height = 1;

    var streaming = false;
    var predicting = false;
    var reloading = false;

    var video = null;
    var canvas = null;
    var resPhoto = null;
    var startbutton = null;
    var userID = Math.floor(Math.random() * 10000000);

    var mapX0 = 0;
    var mapX1 = 100;
    var mapY0 = 0;
    var mapY1 = 100;

    const constraints = {
        audio: false,
        video: {
            facingMode: "environment"
        },
        width: 600,
        height: 600,
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
            resPhoto.removeAttribute("hidden");
            predicting = true;
            startbutton.innerHTML = 'reset';
            findMap();
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

    function findMap() {
        height = video.videoHeight;
        width = video.videoWidth;
        logToServer({
                "name":"finding map",
                "width": width,
                "height": height,
        });

        video.setAttribute('width', width);
        video.setAttribute('height', height);
        canvas.setAttribute('width', width);
        canvas.setAttribute('height', height);

        resPhoto.setAttribute('src', "../static/map-locating.png");
        var context = canvas.getContext('2d');
        if (width && height) {
            canvas.width = width;
            canvas.height = height;
            context.drawImage(video, 0, 0, width, height);

            var canvasData = canvas.toDataURL('image/png');

            $.ajax({
                type: "POST",
                url: "https://ungank.com/findmap",
                data: {
                    userID: userID,
                    imgBase64: canvasData
                },
                async: true,
                headers: {
                    "cache-control": "no-cache"
                }
            }).done(function (d) {
                resPhoto.setAttribute('src', d["minimap"]);
                mapX0 = d["x0"];
                mapX1 = d["x1"];
                mapY0 = d["y0"];
                mapY1 = d["y1"];
                logToServer({
                    "name":"map location received",
                    "x0": d["x0"],
                    "x1": d["x1"],
                    "y0": d["y0"],
                    "y1": d["y1"],
                });

                predictPositions()
            });
        }
    }

    function predictPositions() {
        logToServer({
            "name":"started prediction",
            "width": width,
            "height": height,
        });
        var context = canvas.getContext('2d');
        if (width && height) {
            canvas.width = width;
            canvas.height = height;
            context.drawImage(video, 0, 0, width, height);

            var canvasData = canvas.toDataURL('image/png');
            $.ajax({
                type: "POST",
                url: "https://ungank.com/predict",
                data: {
                    userID: userID,
                    imgBase64: canvasData,
                    x0: mapX0,
                    x1: mapX1,
                    y0: mapY0,
                    y1: mapY1,
                }
            }).done(function (d) {
                logToServer("prediction received");
                resPhoto.setAttribute('src', d["result"]);
                if (predicting) {
                    predictPositions();
                }
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
    canvas = document.getElementById('canvas');
    resPhoto = document.getElementById('response');
    startbutton = document.getElementById('startbutton');

    startbutton.addEventListener('click', handleButtonClick, false);
}