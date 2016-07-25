module.exports = angular.module('trafficPortal.deliveryService.view.chart', [])
    .controller('DeliveryServiceViewChartsController', require('./DeliveryServiceViewChartsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryService.view.chart', {
                url: '/chart',
                abstract: true,
                views: {
                    deliveryServiceViewContent: {
                        templateUrl: 'modules/private/deliveryService/view/charts/deliveryService.view.charts.tpl.html',
                        controller: 'DeliveryServiceViewChartsController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
