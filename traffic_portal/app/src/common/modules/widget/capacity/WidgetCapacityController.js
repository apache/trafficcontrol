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

var WidgetCapacityController = function($scope, $interval, cdnService, propertiesModel) {

	var interval,
		autoRefresh = propertiesModel.properties.dashboard.autoRefresh;

	var getCapacity = function() {
		cdnService.getCapacity()
			.then(function(response) {
				$scope.availablePercent = response.availablePercent;
				$scope.utilizedPercent = response.utilizedPercent;
				$scope.maintenancePercent = response.maintenancePercent;
				$scope.unavailablePercent = response.unavailablePercent;

				var data = [];

				data.push({
					label: "Available",
					color: '#1ABB9C',
					data: $scope.availablePercent
				});
				data.push({
					label: "Utilized",
					color: '#3498DB',
					data: $scope.utilizedPercent
				});
				data.push({
					label: "Maintenance",
					color: '#73879C',
					data: $scope.maintenancePercent
				});
				data.push({
					label: "Down",
					color: '#E74C3C',
					data: $scope.unavailablePercent
				});

				buildGraph(data);
			});

	};

	var buildGraph = function(graphData) {

		var options = {
			series: {
				pie: {
					show: true,
					innerRadius: 0.5,
					radius: 1,
					label: {
						show: false
					}
				}
			},
			grid: {
				hoverable: true
			},
			tooltip: true,
			tooltipOpts: {
				cssClass: "capacityChartTooltip",
				content: "%s: %p.2%",
				defaultTheme: false
			},
			legend: {
				show: false
			}
		};

		$.plot($("#capacityChart"), graphData, options);
	};

	var createInterval = function() {
		killInterval();
		interval = $interval(function() { getCapacity() }, propertiesModel.properties.dashboard.capacity.refreshRateInMS );
	};

	var killInterval = function() {
		if (angular.isDefined(interval)) {
			$interval.cancel(interval);
			interval = undefined;
		}
	};

	$scope.availablePercent = 0;
	$scope.utilizedPercent = 0;
	$scope.maintenancePercent = 0;
	$scope.unavailablePercent = 0;

	$scope.$on("$destroy", function() {
		killInterval();
	});

	var init = function() {
		getCapacity();
		if (autoRefresh) {
			createInterval();
		}
	};
	init();

};

WidgetCapacityController.$inject = ['$scope', '$interval', 'cdnService', 'propertiesModel'];
module.exports = WidgetCapacityController;
