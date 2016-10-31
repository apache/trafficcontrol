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

var ChartTransactionsPerSecondController = function(entity, showSummary, $rootScope, $scope, $uibModal, $q, $timeout, $filter, propertiesModel, dateUtils, statsService) {

    $scope.chartName = propertiesModel.properties.charts.transactionsPerSecond.name;

    var chartDatesChanged = false,
        chartStart,
        chartEnd;

    var summaryStart,
        summaryEnd;

    var chartRangeTimer;

    var loadTransactions = function(start, end) {
        if (!entity || !chartDatesChanged) return;
        chartDatesChanged = false;
        $scope.transactionChartDates = {
            start: start,
            end: end
        };
        getTransactions(start, end);
        $scope.refreshTpsSummaryMetrics(0);
    };

    var getTransactions = function(start, end) {
        var exclude = '',
            ignoreLoadingBar = true,
            showError = true,
            promises = [];

        // edge transactions
        promises.push(statsService.getEdgeTransactions(entity, start, end, $scope.transactionsChartInterval, exclude, ignoreLoadingBar, showError));

        $q.all(promises)
            .then(
            function(responses) {
                // set chart data
                var edgeTransactionsChartData = buildTransactionsChartData(responses[0], start, false);
                $timeout(function () {
                    buildTransactionsChart(edgeTransactionsChartData);
                }, 100);
            },
            function(fault) {
                buildTransactionsChart([]); // build an empty chart
            }).finally(function() {
                $scope.transactionsLoaded = true;
            });
    };

    var buildTransactionsChartData = function(result, start, incremental) {
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

    var buildTransactionsChart = function(edgeTransactionsChartData) {

        var options = {
            xaxis: {
                mode: "time",
                timezone: "browser",
                twelveHourClock: true
            },
            yaxes: [
                {
                    position: "left",
                    axisLabel: "TPS",
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
                    tooltipString += '<span>' + label + ': ' + $filter('number')(yval, 2) + ' TPS</span><br>'
                    return tooltipString;
                }
            }
        };

        $.plot($("#tps-chart"), [ { label: "Edge", data: edgeTransactionsChartData } ], options);

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
        promises.push(statsService.getEdgeTransactionsSummary(entity, start, end, $scope.transactionsChartInterval, exclude, ignoreLoadingBar, showError));

        $q.all(promises)
            .then(
            function(responses) {
                var edgeSummary = responses[0].summary;
                if (angular.isDefined(edgeSummary)) {
                    $scope.tpsEdgeSummary = edgeSummary;
                } else {
                    $scope.resetEdgeSummary();
                }
            },
            function(fault) {
                $scope.resetEdgeSummary();
            }).finally(function() {
                $scope.updatingTpsSummaryMetrics = false;
            });
    };

    var onDateChange = function(args) {
        chartDatesChanged = true;
        chartStart = args.start;
        chartEnd = args.end;
        summaryStart = args.start;
        summaryEnd = args.end;
        updateChartDates(chartStart, chartEnd);
        loadTransactions(chartStart, chartEnd);
    };

    $scope.showSummary = showSummary;

    $scope.updatingTpsSummaryMetrics = false;

    $scope.transactionsLoaded = false;

    $scope.transactionsChartInterval = '60s';

    $scope.resetEdgeSummary = function() {
        $timeout(function() {
            $scope.tpsEdgeSummary = {
                max: 0,
                min: 0,
                totalTransactions: 0,
                average: 0,
                fifthPercentile: 0,
                ninetyFifthPercentile: 0,
                ninetyEighthPercentile: 0
            };
        });
    };
    $scope.resetEdgeSummary();

    $scope.refreshTpsSummaryMetrics = function(delay) {
        if (!$scope.showSummary) return; // don't bother. summary hidden...

        $timeout(function() { $scope.updatingTpsSummaryMetrics = true; });
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
        $scope.refreshTpsSummaryMetrics(0);
    };

    $scope.$on('chartModel::dateChange', function(event, args) {
        onDateChange(args);
    });

    $scope.$on('chartModel::dateRoll', function(event, args) {
        onDateChange(args);
    });

};

ChartTransactionsPerSecondController.$inject = ['entity', 'showSummary', '$rootScope', '$scope', '$uibModal', '$q', '$timeout', '$filter', 'propertiesModel', 'dateUtils', 'statsService'];
module.exports = ChartTransactionsPerSecondController;
