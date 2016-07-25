module.exports = angular.module('trafficPortal.deliveryService', [])
    .controller('DeliveryServiceController', require('./DeliveryServiceController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryService', {
                url: 'delivery-service',
                abstract: true,
                views: {
                    privateContent: {
                        templateUrl: 'modules/private/deliveryService/deliveryService.tpl.html',
                        controller: 'DeliveryServiceController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
