/**
 * @author Ramiro Arenivar
 * For CSCI 5799
 */

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

    // sample videos
    var theVideos = {
        "Videos": [
            { "name": "sample1", "date": "12/12/16", "url": "http://www.w3schools.com/tags/movie.mp4" },
            { "name": "sample2", "date": "12/12/16", "url": "http://www.w3schools.com/tags/movie.mp4" }
        ]
    };

    AvailableVideos.remove({});
    for (var i = 0; i < 2; i++) {
        AvailableVideos.insert(theVideos.Videos[i]);
    }
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