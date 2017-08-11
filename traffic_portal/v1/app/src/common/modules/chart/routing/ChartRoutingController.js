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

var ChartRoutingController = function(entityId, service, $rootScope, $scope, $interval, $filter) {

    var routingInterval,
        routingLoaded = false;

    var getRoutingMethods = function(showTimoutError) {
        if (!$rootScope.online) return;

        var ignoreLoadingBar = true;
        service.getRoutingMethods(entityId, ignoreLoadingBar, showTimoutError)
            .then(
            function(response) {
                routingLoaded = true;
                var staticRoute = Math.round(response.staticRoute * 100) / 100,
                    dsr = Math.round(response.dsr * 100) / 100,
                    err = Math.round(response.err * 100) / 100,
                    miss = Math.round(response.miss * 100) / 100,
                    geo = Math.round(response.geo * 100) / 100,
                    fed = Math.round(response.fed * 100) / 100,
                    cz = Math.round(response.cz * 100) / 100;

                var data = [];

                if (staticRoute > 0) {
                    data.push({
                        label: "Static",
                        color: '#cccccc',
                        data: [ [ staticRoute, 1 ] ]
                    });
                }
                if (dsr > 0) {
                    data.push({
                        label: "DSR",
                        color: '#3c763d',
                        data: [ [ dsr, 1 ] ]
                    });
                }
                if (err > 0) {
                    data.push({
                        label: "Error",
                        color: '#FF0000',
                        data: [ [ err, 1 ] ]
                    });
                }
                if (miss > 0) {
                    data.push({
                        label: "Miss",
                        color: '#a94442',
                        data: [ [ miss, 1 ] ]
                    });
                }
                if (geo > 0) {
                    data.push({
                        label: "3rd Party",
                        color: '#263C53',
                        data: [ [ geo, 1 ] ]
                    });
                }
                if (cz > 0) {
                    data.push({
                        label: "Native",
                        color: '#357EBD',
                        data: [ [ cz, 1 ] ]
                    });
                }
                if (fed > 0) {
                    data.push({
                        label: "Federated",
                        color: '#8a00e6',
                        data: [ [ fed, 1 ] ]
                    });
                }

                buildRoutingChart(data);
            });
    };

    var buildRoutingChart = function(data) {

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

        $.plot($("#routing-chart"), data, options);
    };

    angular.element(document).ready(function () {
        getRoutingMethods(true);
        routingInterval = $interval(function() { getRoutingMethods(false) }, 5 * 60 * 1000); // every 5 mins routing will refresh
    });

    $scope.$on("$destroy", function() {
        if (angular.isDefined(routingInterval)) {
            $interval.cancel(routingInterval);
            routingInterval = undefined;
        }
    });

};

ChartRoutingController.$inject = ['entityId', 'service', '$rootScope', '$scope', '$interval', '$filter'];
module.exports = ChartRoutingController;
