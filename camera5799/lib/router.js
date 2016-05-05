/**
 * @author Ramiro Arenivar
 * For CSCI 5799
 */

Router.configure({
  layoutTemplate: 'layout',
  loadingTemplate: 'loading',
  notFoundTemplate: 'notFound'
});

Router.route('/', {
  name: 'homePage'
});

Router.route('/addCamera/', {
  name: 'addCamera'
});

Router.route('/browseVideos/', {
  name: 'browseVideos'
});

Router.route('/browseCameras/', {
  name: 'browseCameras'
});

Router.route('/editCamera/:name', {
  name: 'editCamera',
  data: function () {
    return {
      cameraToEdit: AvailableCameras.findOne({name: this.params.name})
    };
  }
});

Router.route('/cameraDetail/:name', {
  name: 'cameraDetail',
  data: function() {
    return {
      name: this.params.name
    };
  }
});

Router.route('/editUser/:username', {
  name: 'editUser',
  data:function() {
    return {
      userData: UserData.findOne()
    }
  }
});