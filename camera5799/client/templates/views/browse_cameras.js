Template.browseCameras.helpers({
     availableCameras: function() {
         return AvailableCameras.find();
     }
});

Template.browseCameras.onCreated(function() {
    //alert("on created method fired!");

    var username = Utilities.getUsername();
    var token = Utilities.getUserToken();

    Meteor.call('getCameras', token, username, function(err, res) {
        if (err) {
            alert("Error... " + JSON.stringify(err));
        } else {
            if (res.statusCode == 200) {
                if (res.hasOwnProperty('content')) {
                    res = JSON.parse(res.content);
                    if (res.hasOwnProperty('CameraData')) {
                        AvailableCameras.remove({});
                        for (var i = 0; i < res.CameraData.length; i++) {
                            AvailableCameras.insert(res.CameraData[i]);
                        }
                    }
                }
            }
        }
        return res;
    });

});