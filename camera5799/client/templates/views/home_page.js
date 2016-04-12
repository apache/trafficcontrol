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
              alert(err.message);
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
                  alert(res.content);
              }
          }
           return res;
       });
   }
});