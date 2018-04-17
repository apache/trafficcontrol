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

var TableServerDeliveryServicesController = function(server, serverDeliveryServices, $scope, $state, $uibModal, dateUtils, deliveryServiceUtils, locationUtils, serverUtils, deliveryServiceService, serverService, propertiesModel) {

	var protocols = deliveryServiceUtils.protocols;

	var qstrings = deliveryServiceUtils.qstrings;

	var removeDeliveryService = function(dsId) {
		deliveryServiceService.deleteDeliveryServiceServer(dsId, $scope.server.id)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.server = server;

	$scope.serverDeliveryServices = serverDeliveryServices;

	$scope.showChartsButton = propertiesModel.properties.deliveryServices.charts.show;

	$scope.openCharts = deliveryServiceUtils.openCharts;

	$scope.isEdge = serverUtils.isEdge;

	$scope.protocol = function(ds) {
		return protocols[ds.protocol];
	};

	$scope.qstring = function(ds) {
		return qstrings[ds.qstringIgnore];
	};

	$scope.getRelativeTime = dateUtils.getRelativeTime;

	$scope.cloneDsAssignments = function() {
		var params = {
			title: 'Clone Delivery Service Assignments',
			message: "Please select an edge cache to assign these " + serverDeliveryServices.length + " delivery services to.<br><br>Warning - Any delivery services currently assigned to the selected edge cache will be lost and replaced with these delivery service assignments...",
			labelFunction: function(item) { return item['hostName'] + '.' + item['domainName'] }
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function(serverService) {
					return serverService.getServers({ type: 'EDGE', orderby: 'hostName' });
				}
			}
		});
		modalInstance.result.then(function(selectedServer) {
			var dsIds = _.pluck(serverDeliveryServices, 'id');
			serverService.assignDeliveryServices(selectedServer, dsIds, true, true)
				.then(
					function() {
						locationUtils.navigateToPath('/servers/' + selectedServer.id + '/delivery-services');
					}
				);
		}, function () {
			// do nothing
		});
	};


	$scope.selectDeliveryServices = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/serverDeliveryServices/table.assignDeliveryServices.tpl.html',
			controller: 'TableAssignDeliveryServicesController',
			size: 'lg',
			resolve: {
				server: function() {
					return server;
				},
				deliveryServices: function(deliveryServiceService) {
					return deliveryServiceService.getDeliveryServices({ cdn: server.cdnId });
				},
				assignedDeliveryServices: function() {
					return serverDeliveryServices;
				}
			}
		});
		modalInstance.result.then(function(selectedDsIds) {
			serverService.assignDeliveryServices(server, selectedDsIds, true, false)
				.then(
					function() {
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.confirmRemoveDS = function(ds, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		var params = {
			title: 'Remove Delivery Service from Server?',
			message: 'Are you sure you want to remove ' + ds.xmlId + ' from this server?'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			removeDeliveryService(ds.id);
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
			"columnDefs": [
				{ 'orderable': false, 'targets': 12 }
			],
			"aaSorting": []
		});
	});

};

TableServerDeliveryServicesController.$inject = ['server', 'serverDeliveryServices', '$scope', '$state', '$uibModal', 'dateUtils', 'deliveryServiceUtils', 'locationUtils', 'serverUtils', 'deliveryServiceService', 'serverService', 'propertiesModel'];
module.exports = TableServerDeliveryServicesController;
