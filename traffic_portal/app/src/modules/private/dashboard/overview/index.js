module.exports = angular.module('trafficPortal.private.dashboard.overview', [])
    .controller('DashboardDeliveryServicesController', require('./deliveryServices/DashboardDeliveryServicesController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.dashboard.overview', {
                url: '',
                views: {
                    cacheGroupsContent: {
                        templateUrl: 'common/modules/cacheGroups/cacheGroups.tpl.html',
                        controller: 'CacheGroupsController',
                        resolve: {
                            entityId: function() {
                                return null;
                            },
                            service: function(healthService) {
                                return healthService;
                            },
                            showDownload: function() {
                                return false;
                            }
                        }
                    },
                    deliveryServicesContent: {
                        templateUrl: 'modules/private/dashboard/overview/deliveryServices/dashboard.deliveryServices.tpl.html',
                        controller: 'DashboardDeliveryServicesController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
