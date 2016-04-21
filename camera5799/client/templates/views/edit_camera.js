var editCameraCalls = {

    editCamera: function (editCameraObj) {
        var userName = Utilities.getUsername();
        var token = Utilities.getUserToken();

        Meteor.call('editCameraInformation', token, userName, editCameraObj, function(err, res) {
            if (err) {
                if (err.hasOwnProperty('content')) {
                    alert(JSON.stringify(err.content));
                }
                else {
                    alert(err);
                }
            }
            else {
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('content')) {
                        res = JSON.parse(res.content);
                        if (res.hasOwnProperty('Status') && res.Status == "Success") {
                            if (res.hasOwnProperty('Message')) { alert(res.Message); }
                            Router.go('browseCameras');
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

Template.editCamera.events({

    'click #btn-edit-camera': function (evt, tpl) {
        var name = tpl.find('input#edit-cameraname').value;
        var location = tpl.find('input#edit-cameralocation').value;
        var url = tpl.find('input#edit-cameraurl').value;
        var cameraUsername = tpl.find('input#edit-camerausername').value;
        var cameraPassword = tpl.find('input#edit-camerapassword').value;

        if (name && location && url && cameraUsername && cameraPassword) {
            var cameraObj = {
                name: name,
                location: location,
                url: url,
                username: cameraUsername,
                password: cameraPassword
            };
            editCameraCalls.editCamera(cameraObj);
        }
        else {
            alert("All fields are required!");
        }
    }
});
