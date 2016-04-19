Template.browseCameras.helpers({
     availableCameras: function() {
         return AvailableCameras.find();
     }
});

Template.browseCameras.onCreated(function() {
    //alert("on created method fired!");

    var login_data = Session.get('login_response');
    var token = null;
    if (login_data.hasOwnProperty('token')) {
        token = login_data.token;
    }

    Meteor.call('getCameras', token, function(err, res) {
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