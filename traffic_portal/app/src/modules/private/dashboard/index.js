module.exports = angular.module('trafficPortal.private.dashboard', [])
    .controller('DashboardController', require('./DashboardController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.dashboard', {
                url: 'dashboard',
                abstract: true,
                views: {
                    privateContent: {
                        templateUrl: 'modules/private/dashboard/dashboard.tpl.html',
                        controller: 'DashboardController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
