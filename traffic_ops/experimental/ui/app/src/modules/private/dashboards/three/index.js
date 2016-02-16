module.exports = angular.module('trafficOps.private.dashboards.three', [])
    .controller('DashboardsThreeController', require('./DashboardsThreeController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.dashboards.three', {
                url: '/three',
                views: {
                    dashboardsContent: {
                        templateUrl: 'modules/private/dashboards/three/dashboards.three.tpl.html',
                        controller: 'DashboardsThreeController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
