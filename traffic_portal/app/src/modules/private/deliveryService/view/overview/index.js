module.exports = angular.module('trafficPortal.deliveryService.view.overview', [])
    .controller('DeliveryServiceViewOverviewController', require('./DeliveryServiceViewOverviewController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryService.view.overview', {
                url: '',
                abstract: true,
                views: {
                    deliveryServiceViewContent: {
                        templateUrl: 'modules/private/deliveryService/view/overview/deliveryService.view.overview.tpl.html',
                        controller: 'DeliveryServiceViewOverviewController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
