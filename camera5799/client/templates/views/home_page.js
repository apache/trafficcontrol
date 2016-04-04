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

       Meteor.call('loginCall', username, password, function(err, res) {
           console.log("client response... ", res);
          if (err) {
              console.log("error from client");
              localStorage.setItem('login_response', JSON.stringify({error: err}));
          } else {
              console.log("the response from client! ", res);
              localStorage.setItem('login_response', JSON.stringify({token: res, username: username}));
          }
           Session.set('login_response', JSON.parse(localStorage.getItem('login_response')));
           Router.go('browseCameras');
           return res;

       });
   }
});