var addCameraCalls = {

    registerCamera: function(registerObj) {
        Meteor.call('registerCamera', Utilities.getUserToken(), Utilities.getUsername(), registerObj, function(err, res) {
            console.log("client response res... ", res);
            console.log("client response err... ", err);
            if (err) {
                alert(JSON.stringify(err.content));
            } else {
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('content')) {
                        res = JSON.parse(res.content);
                        if (res.hasOwnProperty('Message')) {
                            res = res.Message;
                            alert(res);
                            Utilities.clearForm('form-register-camera');
                        }
                    }
                }
                else {
                    alert(JSON.stringify(res.content));
                }
            }
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
                URL: cameraURL,
                username: cameraUsername,
                password: cameraPassword
            };
            addCameraCalls.registerCamera(dataObj);
        }
        else {
            alert("Please fill all fields");
        }
    }
});