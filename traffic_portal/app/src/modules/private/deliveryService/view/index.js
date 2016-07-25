module.exports = angular.module('trafficPortal.deliveryService.view', [])
    .controller('DeliveryServiceViewController', require('./DeliveryServiceViewController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryService.view', {
                url: '/{deliveryServiceId}',
                abstract: true,
                views: {
                    deliveryServiceContent: {
                        templateUrl: 'modules/private/deliveryService/view/deliveryService.view.tpl.html',
                        controller: 'DeliveryServiceViewController',
                        resolve: {
                            deliveryService: function(user, deliveryServicesModel, $stateParams) {
                                return deliveryServicesModel.getDeliveryService($stateParams.deliveryServiceId);
                            }
                        }
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
