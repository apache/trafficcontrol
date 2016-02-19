module.exports = angular.module('trafficOps.private.configure.deliveryServices', [])
    .controller('ConfigureDeliveryServicesController', require('./ConfigureDeliveryServicesController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.deliveryServices', {
                url: 'delivery-services',
                views: {
                    configureContent: {
                        templateUrl: 'modules/private/configure/deliveryServices/configure.deliveryServices.tpl.html',
                        controller: 'ConfigureDeliveryServicesController'
                    }
                },
                resolve: {
                    deliveryServices: function(user, deliveryServiceService, deliveryServicesModel) {
                        return deliveryServiceService.getDeliveryServices(false);
                    }
                }

            })
        ;
        $urlRouterProvider.otherwise('/');
    });
