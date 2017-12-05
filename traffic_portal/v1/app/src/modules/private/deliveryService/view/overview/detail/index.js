/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 */

module.exports = angular.module('trafficPortal.deliveryService.view.overview.detail', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryService.view.overview.detail', {
                url: '',
                views: {
                    chartDatesContent: {
                        templateUrl: 'common/modules/chart/dates/chart.dates.tpl.html',
                        controller: 'ChartDatesController',
                        resolve: {
                            customLabel: function() {
                                return 'Delivery Service Bandwidth';
                            },
                            showAutoRefreshBtn: function() {
                                return true;
                            }
                        }
                    },
                    bandwidthContent: {
                        templateUrl: 'common/modules/chart/bandwidthPerSecond/chart.bandwidthPerSecond.tpl.html',
                        controller: 'ChartBandwidthPerSecondController',
                        resolve: {
                            entity: function(user, $stateParams, deliveryServicesModel) {
                                return deliveryServicesModel.getDeliveryService($stateParams.deliveryServiceId);
                            },
                            showSummary: function() {
                                return true;
                            }
                        }
                    },
                    purgeContent: {
                        templateUrl: 'common/modules/tools/purge/tools.purge.tpl.html',
                        controller: 'ToolsPurgeController'
                    },
                    capacityContent: {
                        templateUrl: 'common/modules/chart/capacity/chart.capacity.tpl.html',
                        controller: 'ChartCapacityController',
                        resolve: {
                            entityId: function($stateParams) {
                                return $stateParams.deliveryServiceId;
                            },
                            service: function(deliveryServiceService) {
                                return deliveryServiceService;
                            }
                        }
                    },
                    cacheGroupsContent: {
                        templateUrl: 'common/modules/cacheGroups/cacheGroups.tpl.html',
                        controller: 'CacheGroupsController',
                        resolve: {
                            entityId: function($stateParams) {
                                return $stateParams.deliveryServiceId;
                            },
                            service: function(deliveryServiceService) {
                                return deliveryServiceService;
                            },
                            showDownload: function() {
                                return true;
                            }
                        }
                    },
                    routingContent: {
                        templateUrl: 'common/modules/chart/routing/chart.routing.tpl.html',
                        controller: 'ChartRoutingController',
                        resolve: {
                            entityId: function($stateParams) {
                                return $stateParams.deliveryServiceId;
                            },
                            service: function(deliveryServiceService) {
                                return deliveryServiceService;
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
