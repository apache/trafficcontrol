module.exports = angular.module('trafficOps.private.monitor.dashboards.three', [])
    .controller('DashboardsThreeController', require('./DashboardsThreeController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.monitor.dashboards.three', {
                url: '/three',
                views: {
                    dashboardsContent: {
                        templateUrl: 'modules/private/monitor/dashboards/three/dashboards.three.tpl.html',
                        controller: 'DashboardsThreeController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
