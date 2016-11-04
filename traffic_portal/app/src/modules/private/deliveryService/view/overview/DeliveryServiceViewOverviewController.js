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

var DeliveryServiceViewOverviewController = function($scope, $location, $state, $uibModal, propertiesModel, deliveryServiceService, chartModel) {

    var getFailoverStatus = function() {
        var ignoreLoadingBar = true;
        deliveryServiceService.getState($scope.deliveryService.id, ignoreLoadingBar)
            .then(function(response) {
                $scope.failover = response.failover;
            });
    };

    $scope.properties = propertiesModel.properties;

    $scope.failover = {
        configured: false,
        enabled: false,
        destination: {
            location: null,
            type: ''
        },
        locations: []
    };

    $scope.viewConfig = function(ds) {

        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/deliveryService/config/edit/deliveryService.config.edit.tpl.html',
            controller: 'DSConfigEditController',
            size: 'lg',
            windowClass: 'ds-config-modal',
            resolve: {
                deliveryService: function () {
                    return angular.copy(ds);
                }
            }
        });

        modalInstance.result.then(function() {
        }, function () {
            // do nothing
        });
    };

    $scope.navigateToChart = function(dsId, type) {
        $location.url('/delivery-service/' + dsId + '/chart/' + type).search({ start: moment(chartModel.chart.start).format(), end: moment(chartModel.chart.end).format() });
    };

    angular.element(document).ready(function () {
        if ($scope.deliveryService && $scope.deliveryService.active) {
            getFailoverStatus();
        }
    });

};

DeliveryServiceViewOverviewController.$inject = ['$scope', '$location', '$state', '$uibModal', 'propertiesModel', 'deliveryServiceService', 'chartModel'];
module.exports = DeliveryServiceViewOverviewController;
