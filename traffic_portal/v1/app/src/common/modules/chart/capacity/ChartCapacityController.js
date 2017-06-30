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

var ChartCapacityController = function(entityId, service, $rootScope, $scope, $interval, $filter) {

    var capacityInterval,
        capacityLoaded = false;

    var getCapacity = function(showTimeoutError) {
        if (!$rootScope.online) return;

        var ignoreLoadingBar = true;
        service.getCapacity(entityId, ignoreLoadingBar, showTimeoutError)
            .then(
            function(response) {
                capacityLoaded = true;
                var maintenancePercent = Math.round(response.maintenancePercent * 100) / 100,
                    unavailablePercent = Math.round(response.unavailablePercent * 100) / 100,
                    availablePercent = Math.round(response.availablePercent * 100) / 100,
                    utilizedPercent = Math.round(response.utilizedPercent * 100) / 100;

                var data = [];

                if (maintenancePercent > 0) {
                    data.push({
                        label: "Maintenance",
                        color: '#cccccc',
                        data: [ [ maintenancePercent, 1 ] ]
                    });
                }
                if (unavailablePercent > 0) {
                    data.push({
                        label: "Down",
                        color: '#a94442',
                        data: [ [ unavailablePercent, 1 ] ]
                    });
                }
                if (availablePercent > 0) {
                    data.push({
                        label: "Available",
                        color: '#91ca32',
                        data: [ [ availablePercent, 1 ] ]
                    });
                }
                if (utilizedPercent > 0) {
                    data.push({
                        label: "Utilized",
                        color: '#357ebd',
                        data: [ [ utilizedPercent, 1 ] ]
                    });
                }

                buildCapacityChart(data);
            });
    };

    var buildCapacityChart = function(data) {

        var options = {
            series: {
                stack: true,
                lines: {show: false, steps: false },
                bars: {
                    show: true,
                    horizontal: true,
                    barWidth: 0.9,
                    align: 'center'
                }
            },
            grid: {
                borderWidth: 0,
                hoverable: true
            },
            tooltip: {
                show: true,
                content: function(label, xval, yval, flotItem){
                    return '<span>' + label + ': ' + $filter('number')(xval, 2) + '%</span><br>';
                }
            },
            yaxis: {
                ticks: [[ 1,'%' ]]
            }
        };

        $.plot($("#capacity-chart"), data, options);
    };

    angular.element(document).ready(function () {
        getCapacity(true);
        capacityInterval = $interval(function() { getCapacity(false) }, 5 * 60 * 1000); // every 5 mins capacity will refresh
    });

    $scope.$on("$destroy", function() {
        if (angular.isDefined(capacityInterval)) {
            $interval.cancel(capacityInterval);
            capacityInterval = undefined;
        }
    });

};

ChartCapacityController.$inject = ['entityId', 'service', '$rootScope', '$scope', '$interval', '$filter'];
module.exports = ChartCapacityController;
