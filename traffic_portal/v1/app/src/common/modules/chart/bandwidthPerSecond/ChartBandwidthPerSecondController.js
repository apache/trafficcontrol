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

var ChartBandwidthPerSecondController = function(entity, showSummary, $rootScope, $scope, $uibModal, $q, $timeout, $filter, dateUtils, numberUtils, messageModel, propertiesModel, statsService) {

    $scope.chartName = propertiesModel.properties.charts.bandwidthPerSecond.name;

    var chartDatesChanged = false,
        chartStart,
        chartEnd;

    var summaryStart,
        summaryEnd;

    var chartRangeTimer;

    var loadBandwidth = function(start, end) {
        if (!entity || !chartDatesChanged) return;
        chartDatesChanged = false;
        $scope.bandwidthChartDates = {
            start: start,
            end: end
        };
        getBandwidth(start, end);
        $scope.refreshBpsSummaryMetrics(0);
    };

    var getBandwidth = function(start, end) {

        var exclude = '',
            ignoreLoadingBar = true,
            showError = true,
            promises = [];

        // edge bandwidth
        promises.push(statsService.getEdgeBandwidth(entity, start, end, $scope.bandwidthChartInterval, exclude, ignoreLoadingBar, showError));

        $q.all(promises)
            .then(
                function(responses) {
                    // set chart data
                    var edgeBandwidthChartData = buildEdgeBandwidthChartData(responses[0], start, false);

                    $timeout(function () {
                        buildBandwidthChart(edgeBandwidthChartData);
                    }, 100);
                },
                function(fault) {
                    buildBandwidthChart([]); // build an empty chart
                }).finally(function() {
                    $scope.bandwidthLoaded = true;
                });
    };

    var buildEdgeBandwidthChartData = function(result, start, incremental) {
        var normalizedChartData = [],
            summary = result.summary,
            series = result.series;

        if (angular.isDefined(series)) {
            if (!incremental && angular.isDefined(summary)) {
                $scope.unitSize = numberUtils.shrink(summary.average)[1];
            }
            _.each(series.values, function(seriesItem) {
                if (moment(seriesItem[0]).isSame(start) || moment(seriesItem[0]).isAfter(start)) {
                    if (_.isNumber(seriesItem[1]) || !incremental) {
                        normalizedChartData.push([ moment(seriesItem[0]).valueOf(), numberUtils.convertTo(seriesItem[1], $scope.unitSize) ]); // converts data to appropriate unit
                    }
                }
            });
        }

        return normalizedChartData;
    };

    var buildBandwidthChart = function(edgeBandwidthChartData) {

        var options = {
            xaxis: {
                mode: "time",
                timezone: "browser",
                twelveHourClock: true
            },
            yaxes: [
                {
                    position: "left",
                    axisLabel: $scope.unitSize + "ps",
                    axisLabelUseCanvas: true,
                    axisLabelFontSizePixels: 12,
                    axisLabelFontFamily: 'Verdana, Arial',
                    axisLabelPadding: 3
                }
            ],
            grid: { hoverable: true },
            tooltip: {
                show: true,
                content: function(label, xval, yval, flotItem){
                    var tooltipString = dateUtils.dateFormat(xval, "ddd mmm d yyyy h:MM:ss tt (Z)") + '<br>';
                    tooltipString += '<span>' + label + ': ' + $filter('number')(yval, 2) + ' ' + $scope.unitSize + 'ps</span><br>'
                    return tooltipString;
                }
            }
        };

        $.plot($("#bps-chart"), [ { label: "Edge", data: edgeBandwidthChartData } ], options);

    };

    var updateChartDates = function(start, end) {
        $scope.dateRangeText = dateUtils.dateFormat(start.toDate(), "ddd mmm d yyyy h:MM tt (Z)") + ' to ' + dateUtils.dateFormat(end.toDate(), "ddd mmm d yyyy h:MM tt (Z)");
    };

    var getSummaryMetrics = function(start, end) {

        var exclude = 'series',
            ignoreLoadingBar = true,
            showError = false,
            promises = [];

        // edge summary
        promises.push(statsService.getEdgeBandwidthSummary(entity, start, end, $scope.bandwidthChartInterval, exclude, ignoreLoadingBar, showError));

        $q.all(promises)
            .then(
                function(responses) {
                    var edgeSummary = responses[0].summary;
                    if (angular.isDefined(edgeSummary)) {
                        $scope.bpsEdgeSummary = edgeSummary;
                    } else {
                        $scope.resetEdgeSummary();
                    }
                },
                function(fault) {
                    $scope.resetEdgeSummary();
                }).finally(function() {
                    $scope.updatingBpsSummaryMetrics = false;
                });
    };

    var onDateChange = function(args) {
        chartDatesChanged = true;
        chartStart = args.start;
        chartEnd = args.end;
        summaryStart = args.start;
        summaryEnd = args.end;
        updateChartDates(chartStart, chartEnd);
        loadBandwidth(chartStart, chartEnd);
    };

    $scope.showSummary = showSummary;

    $scope.updatingBpsSummaryMetrics = false;

    $scope.bandwidthLoaded = false;

    $scope.bandwidthChartInterval = '60s';

    $scope.unitSize = 'Kb';

    $scope.ratio = numberUtils.ratio;

    $scope.resetEdgeSummary = function() {
        $timeout(function(){
            $scope.bpsEdgeSummary = {
                max: 0,
                min: 0,
                totalBytes: 0,
                average: 0,
                fifthPercentile: 0,
                ninetyFifthPercentile: 0,
                ninetyEighthPercentile: 0
            };
        });
    };
    $scope.resetEdgeSummary();

    $scope.refreshBpsSummaryMetrics = function(delay) {
        if (!$scope.showSummary) return; // don't bother. summary hidden...

        $timeout(function() { $scope.updatingBpsSummaryMetrics = true; });
        if (chartRangeTimer) {
            $timeout.cancel(chartRangeTimer);
        }
        chartRangeTimer = $timeout(function () {
            getSummaryMetrics(summaryStart, summaryEnd);
        }, delay);
    };

    $scope.hideSummaryMetrics = function() {
        $scope.showSummary = false;
        $scope.resetEdgeSummary();
    };

    $scope.showSummaryMetrics = function() {
        $scope.showSummary = true;
        $scope.refreshBpsSummaryMetrics(0);
    };

    $scope.$on('chartModel::dateChange', function(event, args) {
        onDateChange(args);
    });

    $scope.$on('chartModel::dateRoll', function(event, args) {
        onDateChange(args);
    });

};

ChartBandwidthPerSecondController.$inject = ['entity', 'showSummary', '$rootScope', '$scope', '$uibModal', '$q', '$timeout', '$filter', 'dateUtils', 'numberUtils', 'messageModel', 'propertiesModel', 'statsService'];
module.exports = ChartBandwidthPerSecondController;
