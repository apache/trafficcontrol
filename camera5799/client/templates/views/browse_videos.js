Template.browseVideos.helpers({
    availableVideos: function() {
        return AvailableVideos.find();
    }
});

Template.browseVideos.onCreated(function() {

    // TODO: If there are no cameras, server response with a 500 error

    var username = Utilities.getUsername();
    var token = Utilities.getUserToken();
    var titleAlert = 'error';
    var messageAlert = null;
    var typeAlert = 'error';
    var showAlert = true;

    var theVideos = {
        "Videos": [
            { "name": "testname1", "date": "12/12/16", "url": "http://www.w3schools.com/tags/movie.mp4" },
            { "name": "testname2", "date": "12/13/16", "url": "http://www.w3schools.com/tags/movie.mp4" },
            { "name": "testname3", "date": "12/14/16", "url": "http://www.w3schools.com/tags/movie.mp4" },
            { "name": "testname4", "date": "12/15/16", "url": "http://www.w3schools.com/tags/movie.mp4" }
        ]
    };

    AvailableVideos.remove({});
    for (var i = 0; i < 4; i++) {
        AvailableVideos.insert(theVideos.Videos[i]);
    }



    // Meteor.call('getVideos', token, username, function(err, res) {
    //     if (err) {
    //         messageAlert = JSON.stringify(err);
    //     } else {
    //         messageAlert = JSON.stringify(res);
    //         if (res.statusCode == 200) {
    //             if (res.hasOwnProperty('content')) {
    //                 res = JSON.parse(res.content);
    //                 if (res.hasOwnProperty('CameraData')) {
    //                     showAlert = false;
    //                     AvailableCameras.remove({});
    //                     for (var i = 0; i < res.CameraData.length; i++) {
    //                         AvailableCameras.insert(res.CameraData[i]);
    //                     }
    //                 }
    //             }
    //         }
    //     }
    //     if (showAlert) { swal(titleAlert, messageAlert, typeAlert); }
    //     return res;
    // });

});

Template.browseVideos.events({

    'click .video-link': function(evt, tpl) {
        var videoURL = $(evt.currentTarget).attr('data-video-url');
        var videoName = $(evt.currentTarget).text();
        var videoInfo = {
            url: videoURL,
            name: videoName
        }
        Modal.show('videoModal', videoInfo);
    }
});