/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

var WidgetCDNChartController = function(cdn, $scope, $timeout, $filter, $q, cdnService, cacheStatsService, dateUtils, locationUtils, numberUtils) {

	var chartSeries,
		chartOptions;

	var getCDN = function(id) {
		cdnService.getCDN(id)
			.then(function(result) {
				$scope.cdn = result;
				registerResizeListener();
				getCurrentStats($scope.cdn.name);
				getChartData($scope.cdn.name, moment().subtract(1, 'days'), moment().subtract(10, 'seconds'));
			});
	};

	var getCurrentStats = function(cdnName) {
		cdnService.getCurrentStats()
			.then(function(result) {
				$scope.currentStats = _.find(result.currentStats, function(item) {
					return item.cdn == cdnName;
				});
			});
	};

	var getChartData = function(cdnName, start, end) {
		var promises = [];

		// get cdn bandwidth
		promises.push(cacheStatsService.getBandwidth(cdnName, start, end));

		// get cdn connections
		promises.push(cacheStatsService.getConnections(cdnName, start, end));

		$q.all(promises)
			.then(
				function(responses) {
					// set chart data
					var bandwidthChartData = buildBandwidthChartData(responses[0], start),
						connectionsChartData = buildConnectionsChartData(responses[1], start);

					$timeout(function () {
						buildChart(bandwidthChartData, connectionsChartData);
					}, 100);
				},
				function(fault) {
					buildChart([], []); // build an empty chart
				});

	};

	var buildBandwidthChartData = function(result, start) {
		var normalizedChartData = [],
			series = result.series;

		if (angular.isDefined(series)) {
			_.each(series.values, function(seriesItem) {
				if (moment(seriesItem[0]).isSame(start) || moment(seriesItem[0]).isAfter(start)) {
					normalizedChartData.push([ moment(seriesItem[0]).valueOf(), numberUtils.convertTo(seriesItem[1], $scope.unitSize) ]); // converts data to appropriate unit
				}
			});
		}

		return normalizedChartData;
	};

	var buildConnectionsChartData = function(result, start) {
		var normalizedChartData = [],
			series = result.series;

		if (angular.isDefined(series)) {
			_.each(series.values, function(seriesItem) {
				if (moment(seriesItem[0]).isSame(start) || moment(seriesItem[0]).isAfter(start)) {
					if (_.isNumber(seriesItem[1])) {
						normalizedChartData.push([ moment(seriesItem[0]).valueOf(), seriesItem[1] ]);
					}
				}
			});
		}

		return normalizedChartData;
	};

	var buildChart = function(bandwidthChartData, connectionsChartData) {

		chartOptions = {
			xaxis: {
				mode: "time",
				timezone: "utc",
				twelveHourClock: false
			},
			yaxes: [
				{
					position: "left",
					axisLabel: "Bandwidth (Gbps)",
					axisLabelUseCanvas: true,
					axisLabelFontSizePixels: 12,
					axisLabelFontFamily: 'Verdana, Arial',
					axisLabelPadding: 3
				},
				{
					position: "right",
					axisLabel: "Connections",
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
					var tooltipString = dateUtils.dateFormat(xval, "UTC: ddd mmm d yyyy H:MM:ss tt (Z)") + '<br>';
					tooltipString += '<span>' + label + ': ' + $filter('number')(yval, 2) + '</span><br>'
					return tooltipString;
				}
			}
		};

		chartSeries = [
			{ label: "Bandwidth", yaxis: 1, color: '#3498DB', data: bandwidthChartData },
			{ label: "Connections", yaxis: 2, color: '#E74C3C', data: connectionsChartData }
		];

		plotChart();

	};

	var registerResizeListener = function() {
		$(window).resize(plotChart);
	};

	var plotChart = function() {
		if (chartOptions && chartSeries) {
			$.plot($("#bps-chart-" + $scope.cdn.id), chartSeries, chartOptions);
		}
	};

	$scope.cdn;

	$scope.unitSize = 'Gb';

	$scope.randomId = '_' + Math.random().toString(36).substr(2, 9);

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		var cdnId;
		if (cdn) {
			cdnId = cdn.id;
		} else {
			// cdn wasn't passed in. we need to figure it out on our own
			cdnId = $('#' + $scope.randomId).closest('.chartContainer').data('cdnid');
		}
		getCDN(cdnId);
	});

};

WidgetCDNChartController.$inject = ['cdn', '$scope', '$timeout', '$filter', '$q', 'cdnService', 'cacheStatsService', 'dateUtils', 'locationUtils', 'numberUtils'];
module.exports = WidgetCDNChartController;
