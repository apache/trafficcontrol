module.exports = angular.module('trafficOps.private.configure.deliveryServices.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.deliveryServices.edit', {
                url: '/{deliveryServiceId}/edit',
                views: {
                    deliveryServicesContent: {
                        templateUrl: 'common/modules/form/deliveryService/form.deliveryService.tpl.html',
                        controller: 'FormEditDeliveryServiceController',
                        resolve: {
                            deliveryService: function($stateParams, deliveryServiceService) {
                                return deliveryServiceService.getDeliveryService($stateParams.deliveryServiceId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
