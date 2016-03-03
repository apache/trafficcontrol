module.exports = angular.module('trafficOps.private.monitor.dashboards.two', [])
    .controller('DashboardsTwoController', require('./DashboardsTwoController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.monitor.dashboards.two', {
                url: '/two',
                views: {
                    dashboardsContent: {
                        templateUrl: 'modules/private/monitor/dashboards/two/dashboards.two.tpl.html',
                        controller: 'DashboardsTwoController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
