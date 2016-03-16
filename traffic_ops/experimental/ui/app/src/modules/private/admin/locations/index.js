module.exports = angular.module('trafficOps.private.admin.locations', [])
    .controller('LocationsController', require('./LocationsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.locations', {
                url: '/locations',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/locations/locations.tpl.html',
                        controller: 'LocationsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
