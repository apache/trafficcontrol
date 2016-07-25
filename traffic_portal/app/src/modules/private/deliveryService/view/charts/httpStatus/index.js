module.exports = angular.module('trafficPortal.deliveryService.view.chart.httpStatus', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryService.view.chart.httpStatus', {
                url: '/http-status-by-class',
                views: {
                    chartDatesContent: {
                        templateUrl: 'common/modules/chart/dates/chart.dates.tpl.html',
                        controller: 'ChartDatesController',
                        resolve: {
                            customLabel: function() {
                                return 'Data';
                            },
                            showAutoRefreshBtn: function() {
                                return true;
                            }
                        }
                    },
                    chartContent: {
                        templateUrl: 'common/modules/chart/httpStatus/chart.httpStatus.tpl.html',
                        controller: 'ChartHttpStatusController',
                        resolve: {
                            entity: function(user, $stateParams, deliveryServicesModel) {
                                return deliveryServicesModel.getDeliveryService($stateParams.deliveryServiceId);
                            }
                        }
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
