module.exports = angular.module('trafficOps.private.dashboards.one', [])
    .controller('DashboardsOneController', require('./DashboardsOneController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.dashboards.one', {
                url: '/one',
                views: {
                    dashboardsContent: {
                        templateUrl: 'modules/private/dashboards/one/dashboards.one.tpl.html',
                        controller: 'DashboardsOneController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
