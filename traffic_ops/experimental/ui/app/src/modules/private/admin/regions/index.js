module.exports = angular.module('trafficOps.private.admin.regions', [])
    .controller('RegionsController', require('./RegionsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.regions', {
                url: '/regions',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/regions/regions.tpl.html',
                        controller: 'RegionsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
