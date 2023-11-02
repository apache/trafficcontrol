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

/**
 * @param {*} serverCapabilities
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IWindowService} $window
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/ServerCapabilityService")} serverCapabilityService
 * @param {import("../../../models/MessageModel")} messageModel
 */
var TableServerCapabilitiesController = function(serverCapabilities, $scope, $state, $uibModal, $window, locationUtils, serverCapabilityService, messageModel) {

	var deleteServerCapability = function(serverCapability) {
		serverCapabilityService.deleteServerCapability(serverCapability.name)
			.then(function(result) {
				messageModel.setMessages(result.data.alerts, false);
				$scope.refresh();
			});
	};

	var confirmDelete = function(serverCapability) {
		var params = {
			title: 'Delete Server Capability: ' + serverCapability.name,
			key: serverCapability.name
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
			controller: 'DialogDeleteController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			deleteServerCapability(serverCapability);
		});
	};

	$scope.serverCapabilities = serverCapabilities;

	$scope.contextMenuItems = [
		{
			text: 'Open in New Tab',
			click: function ($itemScope) {
				$window.open('/#!/server-capabilities/' + $itemScope.sc.name, '_blank');
			}
		},
		null, // Dividier
		{
			text: 'Edit',
			click: function ($itemScope) {
				$scope.editServerCapability($itemScope.sc.name);
			}
		},
		{
			text: 'Delete',
			click: function ($itemScope) {
				confirmDelete($itemScope.sc);
			}
		},
		null, // Dividier
		{
			text: 'View Delivery Services',
			click: function ($itemScope) {
				locationUtils.navigateToPath('/server-capabilities/delivery-services?name=' + $itemScope.sc.name);
			}
		},
		{
			text: 'View Servers',
			click: function ($itemScope) {
				locationUtils.navigateToPath('/server-capabilities/servers?name=' + $itemScope.sc.name );
			}
		}
	];

	$scope.createServerCapability = function() {
		locationUtils.navigateToPath('/server-capabilities/new');
	};

	$scope.editServerCapability = function(name) {
		locationUtils.navigateToPath('/server-capabilities/edit?name=' + name);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	angular.element(document).ready(function () {
		$('#serverCapabilitiesTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableServerCapabilitiesController.$inject = ['serverCapabilities', '$scope', '$state', '$uibModal', '$window', 'locationUtils', 'serverCapabilityService', 'messageModel'];
module.exports = TableServerCapabilitiesController;
