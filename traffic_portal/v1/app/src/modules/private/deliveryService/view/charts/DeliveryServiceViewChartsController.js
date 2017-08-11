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

var DeliveryServiceViewChartsController = function($scope, $location, $uibModal, $state, deliveryServicesModel, propertiesModel, chartModel) {

    $scope.deliveryServices = deliveryServicesModel.deliveryServices;

    $scope.properties = propertiesModel.properties;

    $scope.bpsPopover = {
        title: propertiesModel.properties.charts.bandwidthPerSecond.title,
        content: propertiesModel.properties.charts.bandwidthPerSecond.description
    };

    $scope.tpsPopover = {
        title: propertiesModel.properties.charts.transactionsPerSecond.title,
        content: propertiesModel.properties.charts.transactionsPerSecond.description
    };

    $scope.httpPopover = {
        title: propertiesModel.properties.charts.httpStatus.title,
        content: propertiesModel.properties.charts.httpStatus.description
    };

    $scope.isState = function(state) {
        return $state.current.name == state;
    };

    $scope.changeDS = function(dsId) {
        $state.go($state.current.name, { deliveryServiceId: dsId }, { reload: true });
    };

    $scope.navigateToChart = function(dsId, type) {
        $location.url('/delivery-service/' + dsId + '/chart/' + type).search({ start: moment(chartModel.chart.start).format(), end: moment(chartModel.chart.end).format() });
    };

    $scope.navigateToDeliveryService = function(dsId) {
        $location.url('/delivery-service/' + dsId).search({ start: moment(chartModel.chart.start).format(), end: moment(chartModel.chart.end).format() });
    };

};

DeliveryServiceViewChartsController.$inject = ['$scope', '$location', '$uibModal', '$state', 'deliveryServicesModel', 'propertiesModel', 'chartModel'];
module.exports = DeliveryServiceViewChartsController;
