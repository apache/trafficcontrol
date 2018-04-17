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

var TableUserDeliveryServicesController = function(user, userDeliveryServices, useTenancy, $scope, $state, $uibModal, dateUtils, deliveryServiceUtils, locationUtils, userService, propertiesModel) {

	var protocols = deliveryServiceUtils.protocols;

	var qstrings = deliveryServiceUtils.qstrings;

	$scope.user = user;

	$scope.userDeliveryServices = userDeliveryServices;

	$scope.useTenancy = useTenancy[0].value;

	$scope.showChartsButton = propertiesModel.properties.deliveryServices.charts.show;

	$scope.openCharts = deliveryServiceUtils.openCharts;

	$scope.protocol = function(ds) {
		return protocols[ds.protocol];
	};

	$scope.qstring = function(ds) {
		return qstrings[ds.qstringIgnore];
	};

	$scope.getRelativeTime = dateUtils.getRelativeTime;

	$scope.removeDS = function(dsId, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		userService.deleteUserDeliveryService(user.id, dsId)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.selectDSs = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/userDeliveryServices/table.userDSUnassigned.tpl.html',
			controller: 'TableUserDSUnassignedController',
			size: 'lg',
			resolve: {
				user: function() {
					return user;
				},
				deliveryServices: function(deliveryServiceService) {
					return deliveryServiceService.getDeliveryServices();
				},
				userDeliveryServices: function() {
					return userDeliveryServices;
				}
			}
		});
		modalInstance.result.then(function(selectedDSIds) {
			console.log(selectedDSIds);
			userService.assignUserDeliveryServices(user.id, selectedDSIds)
				.then(
					function() {
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#deliveryServicesTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableUserDeliveryServicesController.$inject = ['user', 'userDeliveryServices', 'useTenancy', '$scope', '$state', '$uibModal', 'dateUtils', 'deliveryServiceUtils','locationUtils', 'userService', 'propertiesModel'];
module.exports = TableUserDeliveryServicesController;
