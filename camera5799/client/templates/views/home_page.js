/**
 * @author Ramiro Arenivar
 * For CSCI 5799
 */

var homePageCalls = {

    login: function (username, password) {

        var titleAlert = 'Error';
        var messageAlert = null;
        var typeAlert = 'error';
        var showAlert = true;

        Meteor.call('loginCall', username, password, function(err, res) {
            if (err) {
                messageAlert = JSON.stringify(err);
            } else {
                messageAlert = JSON.stringify(res);
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('data')) {
                        var theData = res.data;
                        if (theData.hasOwnProperty('Token')) {
                            showAlert = false;
                            localStorage.setItem('login_response', JSON.stringify({token: theData.Token, username: username}));
                            Session.set('login_response', JSON.parse(localStorage.getItem('login_response')));
                            homePageCalls.userData();
                            Router.go('browseCameras');
                        }
                    }
                }
            }
            if (showAlert) { swal(titleAlert, messageAlert, typeAlert); }
            return res;
        });
    },

    register: function (username, firstName, lastName, password) {

        var titleAlert = 'Error';
        var messageAlert = null;
        var typeAlert = 'error';

        Meteor.call('registerUser', username, firstName, lastName, password, function(err, res) {
            if (err) {
                messageAlert = JSON.stringify(err);
            } else {
                messageAlert = JSON.stringify(res);
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('content')) {
                        res = JSON.parse(res.content);
                        if (res.hasOwnProperty('Message')) {
                            titleAlert = 'Success'
                            messageAlert = res.Message;
                            typeAlert = 'success';
                            homePageCalls.login(username, password);
                        }
                    }
                }
            }
            swal(titleAlert, messageAlert, typeAlert);
            return res;
        });
    },

    userData: function() {
        if (Utilities.getUsername() && Utilities.getUserToken()) {
            Meteor.call('userInfo', Utilities.getUsername(), Utilities.getUserToken(), function(err, res) {
                if (err) {
                    swal('Error', JSON.stringify(err), 'warning');
                } else {
                    if (res.statusCode == 200) {
                        if (res.hasOwnProperty('content')) {
                            res = JSON.parse(res.content);
                            if (res.hasOwnProperty('UserData')) {
                                res = res.UserData[0];
                                UserData.insert(res);
                            }

                        }
                    }
                    else {
                        swal('Error', JSON.stringify(res.content), 'warning');
                    }
                }
                return res;
            });
        } else {
            swal('Error', 'Error trying to get user info, please logout and login again', 'warning');
        }
    }

};

Meteor.startup(function() {
    Session.set('login_response', JSON.parse(localStorage.getItem('login_response')));
});

Template.homePage.helpers({
   login_response: function() {
       return Session.get('login_response');
   }
});



Template.homePage.events({

    'click #btn-login': function (evt, tpl) {

       var username = tpl.find('input#login-username').value;
       var password = tpl.find('input#login-password').value;
       homePageCalls.login(username, password);
    },

    'click #btn-signup': function (evt, tpl) {

        var username = tpl.find('input#register-username').value;
        var firstName = tpl.find('input#register-firstname').value;
        var lastName = tpl.find('input#register-lastname').value;
        var password = tpl.find('input#register-password').value;

        if (username && firstName && lastName && password) {
            homePageCalls.register(username, firstName, lastName, password);
        }
        else {
            swal('All fields required', 'Please fill all form fields', 'info');
        }
    }
});