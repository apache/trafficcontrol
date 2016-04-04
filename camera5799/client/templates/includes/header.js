Template.header.helpers({
  
  activeRouteClass: function(/* route names */) {
    var args = Array.prototype.slice.call(arguments, 0);
    args.pop();
    
    var active = _.any(args, function(name) {
      return Router.current() && Router.current().route.getName() === name
    });
    
    return active && 'active';
  },
  
  login_response: function() {
    return Session.get('login_response');
  }
});

Template.header.events({
  'click #logout_button': function(e) {
    e.preventDefault();
    localStorage.removeItem('login_response');
    Session.set('login_response', null);
    Router.go('homePage');
  }
});