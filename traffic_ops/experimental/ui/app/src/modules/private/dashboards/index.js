module.exports = angular.module('trafficOps.private.dashboards', [])
    .controller('DashboardsController', require('./DashboardsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.dashboards', {
                url: 'dashboards',
                abstract: true,
                views: {
                    privateContent: {
                        templateUrl: 'modules/private/dashboards/dashboards.tpl.html',
                        controller: 'DashboardsController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
