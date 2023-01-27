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
 * @param {*} serverCapability
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../models/MessageModel")} messageModel
 * @param {import("../../../../api/ServerCapabilityService")} serverCapabilityService
 */
var FormEditServerCapabilityController = function(serverCapability, $scope, $controller, $uibModal, locationUtils, messageModel, serverCapabilityService) {

	// extends the FormServerCapabilityController to inherit common methods
	angular.extend(this, $controller('FormServerCapabilityController', { serverCapability: serverCapability, $scope: $scope }));

	var deleteServerCapability = function(serverCapability) {
		serverCapabilityService.deleteServerCapability(serverCapability.name)
			.then(function() {
				locationUtils.navigateToPath('/server-capabilities');
			});
	};

	$scope.serverCapabilityName = serverCapability.name;

	$scope.settings = {
		isNew: false,
		saveLabel: 'Update'
	};

	$scope.confirmDelete = function(serverCapability) {
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

	$scope.save = function(currentName, serverCapability) {
		serverCapabilityService.updateServerCapability(currentName, serverCapability).
			then(function(result) {
				messageModel.setMessages(result.data.alerts, currentName !== serverCapability.name);
				locationUtils.navigateToPath('/server-capabilities/edit?name=' + result.data.response.name);
			});
	};

};

FormEditServerCapabilityController.$inject = ['serverCapability', '$scope', '$controller', '$uibModal', 'locationUtils', 'messageModel', 'serverCapabilityService'];
module.exports = FormEditServerCapabilityController;
