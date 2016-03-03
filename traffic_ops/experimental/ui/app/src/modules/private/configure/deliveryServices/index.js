module.exports = angular.module('trafficOps.private.configure.deliveryServices', [])
    .controller('DeliveryServicesController', require('./DeliveryServicesController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.deliveryServices', {
                url: '/delivery-services',
                abstract: true,
                views: {
                    configureContent: {
                        templateUrl: 'modules/private/configure/deliveryServices/deliveryServices.tpl.html',
                        controller: 'DeliveryServicesController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
