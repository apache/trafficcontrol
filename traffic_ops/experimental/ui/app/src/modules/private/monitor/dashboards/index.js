module.exports = angular.module('trafficOps.private.monitor.dashboards', [])
    .controller('DashboardsController', require('./DashboardsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.monitor.dashboards', {
                url: '/dashboards',
                abstract: true,
                views: {
                    monitorContent: {
                        templateUrl: 'modules/private/monitor/dashboards/dashboards.tpl.html',
                        controller: 'DashboardsController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
