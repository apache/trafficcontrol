module.exports = angular.module('trafficPortal.deliveryService.view.chart.transactionsPerSecond', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryService.view.chart.transactionsPerSecond', {
                url: '/transactions-per-second',
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
                        templateUrl: 'common/modules/chart/transactionsPerSecond/chart.transactionsPerSecond.tpl.html',
                        controller: 'ChartTransactionsPerSecondController',
                        resolve: {
                            entity: function(user, $stateParams, deliveryServicesModel) {
                                return deliveryServicesModel.getDeliveryService($stateParams.deliveryServiceId);
                            },
                            showSummary: function() {
                                return true;
                            }
                        }
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
