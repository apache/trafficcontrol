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

module.exports = angular.module('trafficPortal.deliveryService.new', [])
    .controller('DeliveryServiceNewController', require('./DeliveryServiceNewController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryService.new', {
                url: '/new',
                views: {
                    deliveryServiceContent: {
                        templateUrl: 'modules/private/deliveryService/new/deliveryService.new.tpl.html',
                        controller: 'DeliveryServiceNewController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
