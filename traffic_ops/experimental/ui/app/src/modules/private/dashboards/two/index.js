module.exports = angular.module('trafficOps.private.dashboards.two', [])
    .controller('DashboardsTwoController', require('./DashboardsTwoController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.dashboards.two', {
                url: '/two',
                views: {
                    dashboardsContent: {
                        templateUrl: 'modules/private/dashboards/two/dashboards.two.tpl.html',
                        controller: 'DashboardsTwoController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
