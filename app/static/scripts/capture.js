function startUngank() {
    var width = 1;
    var height = 1;

    var streaming = false;
    var predicting = false;

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
        video: {
            mediaSource: "screen", // whole screen sharing
            width: {max: '3840'},
            height: {max: '2160'},
            frameRate: {max: '1'}
        }
    };

    function startPredictions(){
        resPhoto.setAttribute('src', "../static/map-permission.png");
        navigator.mediaDevices.getUserMedia(constraints)
            .then(function(stream) {
                video.srcObject = stream;
                video.play();
                setTimeout(function(){
                    findMap();
                }, 1000);
            })
            .catch(function(err) {
                console.log("An error occurred: " + err);
                resPhoto.setAttribute('src', "../static/map-error.png");
            });

        video.addEventListener('canplay', function(){
            if (!streaming) {
                height = video.videoHeight;
                width = video.videoWidth;

                video.setAttribute('width', width);
                video.setAttribute('height', height);
                canvas.setAttribute('width', width);
                canvas.setAttribute('height', height);
                streaming = true;
            }
        }, false);
    }

    function findMap(){
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
                }
            }).done(function(d) {
                console.log(d)
                resPhoto.setAttribute('src', d["minimap"]);
                mapX0 = d["x0"];
                mapX1 = d["x1"];
                mapY0 = d["y0"];
                mapY1 = d["y1"];
                predictPositions()
            });
        }
    }

    function predictPositions() {
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
            }).done(function(d) {
                console.log(d);
                resPhoto.setAttribute('src', d["result"]);
                if (predicting) {
                    predictPositions();
                }
            });
        }
    }

    console.log("starting up")
    video = document.getElementById('video');
    canvas = document.getElementById('canvas');
    resPhoto = document.getElementById('response');
    startbutton = document.getElementById('startbutton');

    startbutton.addEventListener('click', function(ev){
        if (predicting) {
            location.reload();
        } else {
            predicting = true;
            startbutton.innerHTML = 'stop';
            startPredictions();
            ev.preventDefault();
        }
    }, false);
}