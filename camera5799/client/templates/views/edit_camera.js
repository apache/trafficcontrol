var editCameraCalls = {

    editCamera: function (editCameraObj) {
        var userName = Utilities.getUsername();
        var token = Utilities.getUserToken();

        // Meteor.call('editCameraInformation', token, userName, editCameraObj, function(err, res) {
        //     if (err) {
        //         alert(JSON.stringify(err.content));
        //     } else {
        //         if (res.statusCode == 200) {
        //             if (res.hasOwnProperty('data')) {
        //                 var theData = res.data;
        //                 if (theData.hasOwnProperty('Token')) {
        //                     localStorage.setItem('login_response', JSON.stringify({token: theData.Token, username: username}));
        //                     Session.set('login_response', JSON.parse(localStorage.getItem('login_response')));
        //                     Router.go('browseCameras');
        //                 }
        //             }
        //         }
        //         else {
        //             alert(JSON.stringify(res.content));
        //         }
        //     }
        //     return res;
        // });
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
            // TODO: add call here
        } else {
            alert("All fields are required!");
        }
    }
});
