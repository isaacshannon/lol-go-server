(function() {
    var width = 1920;    // We will scale the photo width to this
    var height = 0;     // This will be computed based on the input stream
  
    var streaming = false;
  
    var video = null;
    var canvas = null;
    var photo = null;
    var resPhoto = null;
    var startbutton = null;
    var userID = Math.floor(Math.random() * 10000000);
  
    constraints = {
      video: {
          mediaSource: "screen", // whole screen sharing
          width: {max: '1920'},
          height: {max: '1080'},
          frameRate: {max: '1'}
        }
      };
  
    function startup() {
      console.log("starting up")
      video = document.getElementById('video');
      canvas = document.getElementById('canvas');
      resPhoto = document.getElementById('response');
      startbutton = document.getElementById('startbutton');
  
      navigator.mediaDevices.getUserMedia(constraints)
      .then(function(stream) {
          video.srcObject = stream;
          video.play();
      })
      .catch(function(err) {
          console.log("An error occurred: " + err);
      });
  
      video.addEventListener('canplay', function(ev){
        if (!streaming) {
          height = video.videoHeight / (video.videoWidth/width);
  
          video.setAttribute('width', width);
          video.setAttribute('height', height);
          canvas.setAttribute('width', width);
          canvas.setAttribute('height', height);
          streaming = true;
        }
      }, false);
  
      startbutton.addEventListener('click', function(ev){
        takepicture();
        ev.preventDefault();
      }, false);
  
      clearphoto();
    }
  
    function clearphoto() {
      var context = canvas.getContext('2d');
      context.fillStyle = "#AAA";
      context.fillRect(0, 0, canvas.width, canvas.height);
  
      var data = canvas.toDataURL('image/png');
      photo.setAttribute('src', data);
    }
  
    function takepicture() {
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
               imgBase64: canvasData,
               user: userID,
            }
          }).done(function(d) {
              console.log(d)
              resPhoto.setAttribute('src', d["result"]);
              takepicture()
          });
      } else {
        clearphoto();
      }
    }
  
    startup();
  })();
  
  