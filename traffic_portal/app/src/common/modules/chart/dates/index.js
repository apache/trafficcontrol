module.exports = angular.module('trafficPortal.chart.dates', [])
    .controller('ChartDatesController', require('./ChartDatesController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.chart.dates', {
                reloadOnSearch: false
            });
        $urlRouterProvider.otherwise('/');
    });
