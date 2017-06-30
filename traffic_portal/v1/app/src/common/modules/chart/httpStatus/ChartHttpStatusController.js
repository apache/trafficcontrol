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

var ChartHttpStatusController = function(entity, $window, $rootScope, $scope, $uibModal, $q, $timeout, $filter, propertiesModel, dateUtils, statsService) {

    $scope.chartName = propertiesModel.properties.charts.httpStatus.name;

    var chartDatesChanged = false,
        chartStart,
        chartEnd;

    var loadAggregateHttpStatusData = function(start, end) {
        if (!entity || !chartDatesChanged) return;
        chartDatesChanged = false;
        getAggregateHttpStatusData(start, end);
    };

    var getAggregateHttpStatusData = function(start, end) {

        var exclude = 'summary',
            ignoreLoadingBar = false,
            showError = false,
            promises = [];

        promises.push(statsService.getEdgeTransactionsByStatusGroup(entity, '2xx', start, end, $scope.httpStatusChartInterval, exclude, ignoreLoadingBar, showError));
        promises.push(statsService.getEdgeTransactionsByStatusGroup(entity, '3xx', start, end, $scope.httpStatusChartInterval, exclude, ignoreLoadingBar, showError));
        promises.push(statsService.getEdgeTransactionsByStatusGroup(entity, '4xx', start, end, $scope.httpStatusChartInterval, exclude, ignoreLoadingBar, showError));
        promises.push(statsService.getEdgeTransactionsByStatusGroup(entity, '5xx', start, end, $scope.httpStatusChartInterval, exclude, ignoreLoadingBar, showError));

        $q.all(promises)
            .then(
            function(responses) {
                var status2xxChartData = buildHttpStatusChartData(responses[0], start, false),
                    status3xxChartData = buildHttpStatusChartData(responses[1], start, false),
                    status4xxChartData = buildHttpStatusChartData(responses[2], start, false),
                    status5xxChartData = buildHttpStatusChartData(responses[3], start, false);
                $timeout(function () {
                    buildHttpChart(status2xxChartData, status3xxChartData, status4xxChartData, status5xxChartData);
                }, 100);
            },
            function(fault) {
                buildHttpChart([], [], [], []);
            }).finally(function() {
                $scope.httpStatusDataLoaded = true;
            });
    };

    var updateChartDates = function(start, end) {
        $scope.dateRangeText = dateUtils.dateFormat(start.toDate(), "ddd mmm d yyyy h:MM tt (Z)") + ' to ' + dateUtils.dateFormat(end.toDate(), "ddd mmm d yyyy h:MM tt (Z)");
    };

    var buildHttpStatusChartData = function(result, start, incremental) {
        var normalizedChartData = [],
            series = result.series;

        if (angular.isDefined(series)) {
            _.each(series.values, function(seriesItem) {
                if (moment(seriesItem[0]).isSame(start) || moment(seriesItem[0]).isAfter(start)) {
                    if (_.isNumber(seriesItem[1]) || !incremental) {
                        normalizedChartData.push([ moment(seriesItem[0]).valueOf(), seriesItem[1] ]);
                    }
                }
            });
        }

        return normalizedChartData;
    };

    var buildHttpChart = function(status2xxChartData, status3xxChartData, status4xxChartData, status5xxChartData) {

        var options = {
            xaxis: {
                mode: "time",
                timezone: "browser",
                twelveHourClock: true
            },
            yaxes: [
                {
                    position: "left",
                    axisLabel: "Success (2xx and 3xx)",
                    axisLabelUseCanvas: true,
                    axisLabelFontSizePixels: 12,
                    axisLabelFontFamily: 'Verdana, Arial',
                    axisLabelPadding: 3
                },
                {
                    position: "right",
                    axisLabel: "Client Error (4xx)",
                    axisLabelUseCanvas: true,
                    axisLabelFontSizePixels: 12,
                    axisLabelFontFamily: 'Verdana, Arial',
                    axisLabelPadding: 3
                },
                {
                    position: "right",
                    axisLabel: "Server Error (5xx)",
                    axisLabelUseCanvas: true,
                    axisLabelFontSizePixels: 12,
                    axisLabelFontFamily: 'Verdana, Arial',
                    axisLabelPadding: 3
                }
            ],
            grid: {
                hoverable: true,
                axisMargin: 20
            },
            tooltip: {
                show: true,
                content: function(label, xval, yval, flotItem){
                    var tooltipString = dateUtils.dateFormat(xval, "ddd mmm d yyyy h:MM:ss tt (Z)") + '<br>';
                    tooltipString += '<span>' + label + ': ' + $filter('number')(yval, 2) + ' TPS</span><br>'
                    return tooltipString;
                }
            }
        };

        var series = [
            { label: "2xx", yaxis: 1, color: "#91ca32", data: status2xxChartData },
            { label: "3xx", yaxis: 1, color: "#5897fb", data: status3xxChartData },
            { label: "4xx", yaxis: 2, color: "#6859a3", data: status4xxChartData },
            { label: "5xx", yaxis: 3, color: "#a94442", data: status5xxChartData }
        ];

        $.plot($("#http-chart"), series, options);

    };

    var onDateChange = function(args) {
        chartDatesChanged = true;
        chartStart = args.start;
        chartEnd = args.end;
        updateChartDates(chartStart, chartEnd);
        loadAggregateHttpStatusData(chartStart, chartEnd);
    };

    $scope.httpStatusDataLoaded = false;

    $scope.httpStatusChartInterval = '60s';

    $scope.resetStatusCodes = function() {
        $timeout(function(){
            $scope.http2xxCodes = [];
            $scope.http3xxCodes = [];
            $scope.http4xxCodes = [];
            $scope.http5xxCodes = [];
        });
    };
    $scope.resetStatusCodes();

    $scope.$on('chartModel::dateChange', function(event, args) {
        onDateChange(args);
    });

    $scope.$on('chartModel::dateRoll', function(event, args) {
        onDateChange(args);
    });

};

ChartHttpStatusController.$inject = ['entity', '$window', '$rootScope', '$scope', '$uibModal', '$q', '$timeout', '$filter', 'propertiesModel', 'dateUtils', 'statsService'];
module.exports = ChartHttpStatusController;
