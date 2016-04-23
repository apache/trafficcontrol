var editUserPageCalls = {
    
    editUser: function (userObj) {
        var username = Utilities.getUsername();
        var token = Utilities.getUserToken();

        Meteor.call('editUserInfo', username, token, userObj, function(err, res) {
            if (err) {
                alert(JSON.stringify(err));
            } else {
                alert(JSON.stringify(res));
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('content')) {
                        res = JSON.parse(res.content);
                        if (res.hasOwnProperty('Status')) {
                            alert(res.Message);
                            UserData.remove({});
                            userObj.username = username;
                            UserData.insert(userObj);
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
    },

    deleteUser: function () {
        var username = Utilities.getUsername();
        var token = Utilities.getUserToken();

        Meteor.call('deleteUser', username, token, function(err, res) {
            if (err) {
                alert(JSON.stringify(err));
            } else {
                alert(JSON.stringify(res));
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('content')) {
                        res = JSON.parse(res.content);
                        if (res.hasOwnProperty('Status')) {
                            alert(res.Message);
                            localStorage.removeItem('login_response');
                            Session.set('login_response', null);
                            // remove all of the client collections on logout
                            var globalObject=Meteor.isClient?window:global;
                            for(var property in globalObject){
                                var object=globalObject[property];
                                if(object instanceof Meteor.Collection){
                                    object.remove({});
                                }
                            }
                            Router.go('homePage');
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

Template.editUser.events({

    'click #btn-edit-user': function (evt, tpl) {

        var username = tpl.find('input#edit-username-current').value;
        var firstName = tpl.find('input#edit-firstName').value;
        var lastName = tpl.find('input#edit-lastName').value;
        var password = tpl.find('input#edit-password').value;

        if (username && firstName && lastName && password) {
            var userObj = {
                firstName: firstName,
                lastName: lastName,
                password: password
            }
            editUserPageCalls.editUser(userObj);
        } else {
            swal('All fields are required!', 'Please fill all form fields', 'warning');
        }
    },

    'click #btn-delete-user': function(evt, tpl) {
        swal({
            title: "Are you sure?",
            text: "You are about to delete your account",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, delete it!",
            closeOnConfirm: false,
            html: false
        }, function(){
            editUserPageCalls.deleteUser();
        });
    }
});