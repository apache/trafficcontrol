module.exports = angular.module('trafficOps.private.configure.locations', [])
    .controller('LocationsController', require('./LocationsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.locations', {
                url: '/locations',
                abstract: true,
                views: {
                    configureContent: {
                        templateUrl: 'modules/private/configure/locations/locations.tpl.html',
                        controller: 'LocationsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
