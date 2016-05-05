/**
 * @author Ramiro Arenivar
 * For CSCI 5799
 */

Template.browseCameras.helpers({
     availableCameras: function() {
         return AvailableCameras.find();
     }
});

Template.browseCameras.onCreated(function() {

    // TODO: If there are no cameras, server response with a 500 error

    var username = Utilities.getUsername();
    var token = Utilities.getUserToken();
    var titleAlert = 'error';
    var messageAlert = null;
    var typeAlert = 'error';
    var showAlert = true;

    Meteor.call('getCameras', token, username, function(err, res) {
        if (err) {
            messageAlert = JSON.stringify(err);
        } else {
            messageAlert = JSON.stringify(res);
            if (res.statusCode == 200) {
                if (res.hasOwnProperty('content')) {
                    res = JSON.parse(res.content);
                    if (res.hasOwnProperty('CameraData')) {
                        showAlert = false;
                        AvailableCameras.remove({});
                        for (var i = 0; i < res.CameraData.length; i++) {
                            AvailableCameras.insert(res.CameraData[i]);
                        }
                    }
                }
            }
        }
        if (showAlert && res) { swal(titleAlert, messageAlert, typeAlert); }
        return res;
    });

});