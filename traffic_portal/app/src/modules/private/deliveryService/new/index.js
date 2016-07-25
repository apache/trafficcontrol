module.exports = angular.module('trafficPortal.deliveryService.new', [])
    .controller('DeliveryServiceNewController', require('./DeliveryServiceNewController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryService.new', {
                url: '/new',
                views: {
                    deliveryServiceContent: {
                        templateUrl: 'modules/private/deliveryService/new/deliveryService.new.tpl.html',
                        controller: 'DeliveryServiceNewController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
