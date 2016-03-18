module.exports = angular.module('trafficOps.private.configure.deliveryServices.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.deliveryServices.new', {
                url: '/new',
                views: {
                    deliveryServicesContent: {
                        templateUrl: 'common/modules/form/deliveryService/form.deliveryService.tpl.html',
                        controller: 'FormNewDeliveryServiceController',
                        resolve: {
                            deliveryService: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
