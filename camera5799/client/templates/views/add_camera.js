var addCameraCalls = {

    registerCamera: function(registerObj) {

        var titleAlert = 'Error';
        var messageAlert = null;
        var typeAlert = 'error';

        Meteor.call('registerCamera', Utilities.getUserToken(), Utilities.getUsername(), registerObj, function(err, res) {
            if (err) {
                messageAlert = JSON.stringify(err);
            } else {
                messageAlert = JSON.stringify(res);
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('content')) {
                        res = JSON.parse(res.content);
                        if (res.hasOwnProperty('Message')) {
                            Utilities.clearForm('form-register-camera');
                            titleAlert = 'Success';
                            messageAlert = res.Message;
                            typeAlert = 'success';
                            Router.go('browseCameras');
                        }
                    }
                }
            }
            swal(titleAlert, messageAlert, typeAlert);
            return res;
        });
    }

};

Template.addCamera.events({

    'click #btn-register-camera': function (evt, tpl) {

        var cameraName = tpl.find('input#register-cameraname').value;
        var cameraLocation = tpl.find('input#register-cameralocation').value;
        var cameraURL = tpl.find('input#register-cameraurl').value;
        var cameraUsername = tpl.find('input#register-camerausername').value;
        var cameraPassword = tpl.find('input#register-camerapassword').value;

        if (cameraName && cameraLocation && cameraURL && cameraUsername && cameraPassword) {
            var dataObj = {
                name: cameraName,
                location: cameraLocation,
                url: cameraURL,
                username: cameraUsername,
                password: cameraPassword
            };
            addCameraCalls.registerCamera(dataObj);
        }
        else {
            swal('All fields required', 'Please fill all form fields', 'info');
        }
    }
});