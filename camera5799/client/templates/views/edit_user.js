/**
 * @author Ramiro Arenivar
 * For CSCI 5799
 */

var editUserPageCalls = {
    
    editUser: function (userObj) {

        var username = Utilities.getUsername();
        var token = Utilities.getUserToken();

        var titleAlert = 'Error';
        var messageAlert = null;
        var typeAlert = 'error';

        Meteor.call('editUserInfo', username, token, userObj, function(err, res) {
            if (err) {
                messageAlert = JSON.stringify(err);
            } else {
                messageAlert = JSON.stringify(res);
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('content')) {
                        res = JSON.parse(res.content);
                        if (res.hasOwnProperty('Status')) {
                            titleAlert = 'Success';
                            messageAlert = res.Message;
                            typeAlert = 'success';
                            UserData.remove({});
                            userObj.username = username;
                            UserData.insert(userObj);
                            Router.go('browseCameras');
                        }
                    }
                }
            }
            swal(titleAlert, messageAlert, typeAlert);
            return res;
        });
    },

    deleteUser: function () {

        var username = Utilities.getUsername();
        var token = Utilities.getUserToken();

        var titleAlert = 'Error';
        var messageAlert = null;
        var typeAlert = 'error';

        Meteor.call('deleteUser', username, token, function(err, res) {
            if (err) {
                messageAlert = JSON.stringify(err);
            } else {
                messageAlert = JSON.stringify(res);
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('content')) {
                        res = JSON.parse(res.content);
                        if (res.hasOwnProperty('Status')) {
                            titleAlert = 'Success';
                            messageAlert = res.Message;
                            typeAlert = 'success';
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
                swal(titleAlert, messageAlert, typeAlert);
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