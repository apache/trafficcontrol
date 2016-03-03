module.exports = angular.module('trafficOps.private.configure.deliveryServices.edit', [])
    .controller('DeliveryServicesEditController', require('./DeliveryServicesEditController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.deliveryServices.edit', {
                url: '/{dsId}',
                views: {
                    deliveryServicesContent: {
                        templateUrl: 'modules/private/configure/deliveryServices/edit/deliveryServices.edit.tpl.html',
                        controller: 'DeliveryServicesEditController'
                    }
                },
                resolve: {
                    deliveryService: function($stateParams, deliveryServiceService) {
                        return deliveryServiceService.getDeliveryService($stateParams.dsId, false);
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
