module.exports = angular.module('trafficOps.private.configure.deliveryServices.list', [])
    .controller('DeliveryServicesListController', require('./DeliveryServicesListController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.deliveryServices.list', {
                url: '',
                views: {
                    deliveryServicesContent: {
                        templateUrl: 'modules/private/configure/deliveryServices/list/deliveryServices.list.tpl.html',
                        controller: 'DeliveryServicesListController'
                    }
                },
                resolve: {
                    deliveryServices: function(user, deliveryServiceService) {
                        return deliveryServiceService.getDeliveryServices(false);
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
