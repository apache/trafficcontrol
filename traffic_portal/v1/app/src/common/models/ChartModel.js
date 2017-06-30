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

var ChartModel = function($rootScope, $location, $interval, messageModel) {

    var chart,
        model = this;

    var autoRefreshInterval;

    var increment = 1,
        unit = 'm';

    var createAutoRefreshInterval = function() {
        killAutoRefreshInterval();
        autoRefreshInterval = $interval(function() { rollDates() }, (increment*60*1000)); // every X minutes
    };

    var killAutoRefreshInterval = function() {
        if (angular.isDefined(autoRefreshInterval)) {
            $interval.cancel(autoRefreshInterval);
            autoRefreshInterval = undefined;
        }
    };

    var rollDates = function() {
        var locationStart = $location.search().start,
            locationEnd = $location.search().end;

        model.chart.start = moment(model.chart.start).add(increment, unit);
        model.chart.end = moment(model.chart.end).add(increment, unit);

        if (locationStart) {
            $location.search('start', model.chart.start.format());
        }

        if (locationEnd) {
            $location.search('end', model.chart.end.format());
        }

        $rootScope.$broadcast('chartModel::dateRoll', { start: model.chart.start, end: model.chart.end });
    };

    this.resetChart = function() {
        var start = $location.search().start,
            end = $location.search().end;

        chart = {};
        chart.start = moment().subtract(1, 'd');
        chart.end = moment();
        chart.autoRefresh = true;

        if (angular.isDefined(start) && angular.isDefined(end)) {
            if (moment(start).isValid() && moment(end).isValid()) {
                chart.start = moment(start);
                chart.end =  moment(end);
                chart.autoRefresh = false;
            } else {
                messageModel.setMessages([ { level: 'error', text: 'Invalid date format detected. Reverting to default.' } ], true);
            }
        }

        if (chart.autoRefresh) {
            createAutoRefreshInterval();
        } else {
            killAutoRefreshInterval();
        }

        var rangeParams = model.calculateRange(chart.start, chart.end);
        chart.range = chart.end.diff(chart.start, rangeParams.interval) + rangeParams.abbrev;
        if (!rangeParams.exact) {
            chart.range = '~' + chart.range;
        }

        this.chart = chart;
    };

    this.changeDates = function(start, end) {
        $location.search('start', start.format());
        $location.search('end', end.format());
        $rootScope.$broadcast('chartModel::dateChange', { start: start, end: end });
    };

    this.calculateRange = function(start, end) {
        // if greater than 1d, use day, if greater than 1hr, use hr, else minute
        var rangeParams = {};
        if (end.diff(start, 'days', true) >= 1) {
            rangeParams = { interval: 'days', abbrev: 'd', exact: (end.diff(start, 'days', true) % 1) == 0 };
        } else if (end.diff(start, 'hours', true) >= 1) {
            rangeParams = { interval: 'hours', abbrev: 'h', exact: (end.diff(start, 'hours', true) % 1) == 0 };
        } else {
            rangeParams = { interval: 'minutes', abbrev: 'm', exact: (end.diff(start, 'minutes', true) % 1) == 0 };
        }
        return rangeParams;
    };

    $rootScope.$watch(
        function() { return model.chart.autoRefresh; },
        function(newValue, oldValue) {
            if (newValue !== oldValue) {
                if (newValue) {
                    createAutoRefreshInterval();
                } else {
                    killAutoRefreshInterval();
                }
            }
        }
    );

    $rootScope.$watch('online', function(newStatus) {
        if (newStatus === false) {
            model.chart.autoRefresh = false;
        }
    });

    var init = function () {
        model.resetChart();
    };
    init();

};

ChartModel.$inject = ['$rootScope', '$location', '$interval', 'messageModel'];
module.exports = ChartModel;