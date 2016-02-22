module.exports = angular.module('trafficOps.private.monitor.dashboards.one', [])
    .controller('DashboardsOneController', require('./DashboardsOneController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.monitor.dashboards.one', {
                url: '/one',
                views: {
                    dashboardsContent: {
                        templateUrl: 'modules/private/monitor/dashboards/one/dashboards.one.tpl.html',
                        controller: 'DashboardsOneController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
