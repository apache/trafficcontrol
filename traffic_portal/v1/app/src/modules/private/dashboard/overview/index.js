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
