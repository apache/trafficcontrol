Router.configure({
  layoutTemplate: 'layout',
  loadingTemplate: 'loading',
  notFoundTemplate: 'notFound'
});

Router.route('/', {
  name: 'homePage'
});

Router.route('/addCamera', {
  name: 'addCamera'
});

Router.route('/browseVideos', {
  name: 'browseVideos'
});

Router.route('/browseCameras', {
  name: 'browseCameras'
});

Router.route('/cameraDetail', {
  name: 'cameraDetail'
});