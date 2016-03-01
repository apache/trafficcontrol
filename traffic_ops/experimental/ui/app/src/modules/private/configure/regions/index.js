module.exports = angular.module('trafficOps.private.configure.regions', [])
    .controller('RegionsController', require('./RegionsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.regions', {
                url: '/regions',
                abstract: true,
                views: {
                    configureContent: {
                        templateUrl: 'modules/private/configure/regions/regions.tpl.html',
                        controller: 'RegionsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
