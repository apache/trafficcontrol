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

var ChartDatesController = function(customLabel, showAutoRefreshBtn, $scope, $location, $timeout, chartModel, messageModel) {

    var setRange = function() {
        var rangeParams = chartModel.calculateRange($scope.chartData.start, $scope.chartData.end);
        $scope.tempChartData.range = $scope.chartData.end.diff($scope.chartData.start, rangeParams.interval) + rangeParams.abbrev;
        if (!rangeParams.exact) {
            $scope.tempChartData.range = '~' + $scope.tempChartData.range;
        }
        $scope.chartData.range = $scope.tempChartData.range;
    };

    var createNowBtn = function() {
        var $nowBtn = $('<button type="button" class="dates-now-btn btn btn-block action-btn">Now</button>');
        $nowBtn.click(function () {
            angular.element(document.getElementById('rangeInput')).scope().setEndToNow();
        });
        $nowBtn.appendTo($('.end-dropdown .datetimepicker'));
    };

    $scope.chartData = chartModel.chart;

    $scope.tempChartData = {
        range: $scope.chartData.range
    };

    $scope.customLabel = customLabel;

    $scope.showAutoRefreshBtn = showAutoRefreshBtn;

    $scope.closeStart = function() {
        $scope.startDropdown = {
            isopen: false
        };
    };
    $scope.closeStart();

    $scope.closeEnd = function() {
        $scope.endDropdown = {
            isopen: false
        };
    };
    $scope.closeEnd();

    $scope.setStart = function(newDate, oldDate) {
        if (moment(newDate).isAfter()) {
            $scope.chartData.start = moment(oldDate);
            messageModel.setMessages([ { level: 'error', text: "Can't set start date to the future." } ], false);
        } else {
            $scope.chartData.start = moment(newDate);
            setRange();
        }
    };

    $scope.setEnd = function(newDate, oldDate) {
        if (moment(newDate).isAfter()) {
            $scope.chartData.end = moment(oldDate);
            messageModel.setMessages([ { level: 'error', text: "Can't set end date to the future." } ], false);
        } else {
            $scope.chartData.end = moment(newDate);
            setRange();
        }
    };

    $scope.setEndToNow = function() {
        $scope.chartData.end = moment();
        setRange();
    };

    $scope.toggleAutoRefresh = function() {
        $scope.chartData.autoRefresh = !$scope.chartData.autoRefresh;
        if ($scope.chartData.autoRefresh) {
            $scope.applyRange(); // applying the range moves it to current
        }
    };

    $scope.revertRange = function() {
        $scope.tempChartData.range = $scope.chartData.range;
    };

    $scope.applyRange = function() {

        var regex = /(\d+)([h|d|w|m|M]$)/, // range must be in the format 1m, 1h, 2d, 3w, 4M
            params = $scope.tempChartData.range.match(regex);

        if (params && params.length == 3) {
            $scope.chartData.start = moment().subtract(params[1], params[2]);
            $scope.chartData.end = moment();
            $scope.chartData.range = $scope.tempChartData.range;
            $scope.changeDates($scope.chartData.start, $scope.chartData.end);
            $scope.chartData.autoRefresh = true && $scope.showAutoRefreshBtn; // showAutoRefreshBtn trumps all. if no show, no autorefresh...EVER!
        } else {
            messageModel.setMessages([ { level: 'error', text: "Invalid date range. Valid increments are 'm' (minute), 'h' (hour), 'd' (day), 'w' (week) or 'M' (month). Example: '30m', '12h', '3d', '3w', '3M'" } ], false);
        }

        $timeout(function () {
            $('#rangeInput').blur(); // need to blur input to hide popover and for some reason a delay helps
        }, 500);

    };

    $scope.changeDates = function(start, end) {
        if (!start.isValid() || !end.isValid()) {
            messageModel.setMessages([ { level: 'error', text: 'Invalid date format detected. Please fix.' } ], false);
        } else {
            chartModel.changeDates(start, end);
        }
    };

    angular.element(document).ready(function () {
        $scope.changeDates(chartModel.chart.start, chartModel.chart.end);
        createNowBtn();
    });

    var init = function () {
        $scope.chartData.autoRefresh = $scope.showAutoRefreshBtn && $scope.chartData.autoRefresh;
    };
    init();

};

ChartDatesController.$inject = ['customLabel', 'showAutoRefreshBtn', '$scope', '$location', '$timeout', 'chartModel', 'messageModel'];
module.exports = ChartDatesController;
