var homePageCalls = {

    login: function (username, password) {
        Meteor.call('loginCall', username, password, function(err, res) {
            if (err) {
                alert(JSON.stringify(err.content));
            } else {
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('data')) {
                        var theData = res.data;
                        if (theData.hasOwnProperty('Token')) {
                            localStorage.setItem('login_response', JSON.stringify({token: theData.Token, username: username}));
                            Session.set('login_response', JSON.parse(localStorage.getItem('login_response')));
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

    register: function (username, firstName, lastName, password) {
        Meteor.call('registerUser', username, firstName, lastName, password, function(err, res) {
            if (err) {
                alert("Error trying to register... " + JSON.stringify(err));
            } else {
                if (res.statusCode == 200) {
                    if (res.hasOwnProperty('content')) {
                        var theData = JSON.parse(res.content);
                        if (theData.hasOwnProperty('Message')) {
                            alert(theData.Message);
                            homeCallbackFuncs.login(username, password);
                        }
                    }
                }
                else {
                    alert("Status code: " + res.statusCode + ", Response: " + JSON.stringify(res));
                }
            }
            return res;
        });
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
            alert("Please fill all fields");
        }
    }
});