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

var FormEditCapabilityController = function(capability, $scope, $controller, $uibModal, $anchorScroll, $location, locationUtils, capabilityService, messageModel) {

	// extends the FormCapabilityController to inherit common methods
	angular.extend(this, $controller('FormCapabilityController', { capability: capability, $scope: $scope }));

	var deleteCapability = function(cap) {
		capabilityService.deleteCapability(cap)
			.then(function(result) {
				messageModel.setMessages(result.alerts, true);
				locationUtils.navigateToPath('/capabilities');
			});
	};

	var save = function(cap) {
		capabilityService.updateCapability(cap).
			then(function(result) {
				$scope.capName = angular.copy(cap.name);
				messageModel.setMessages(result.alerts, false);
				$anchorScroll(); // scrolls window to top
			});
	};

	$scope.capName = angular.copy($scope.capability.name);

	$scope.settings = {
		isNew: false,
		saveLabel: 'Update'
	};

	$scope.viewEndpoints = function() {
		$location.path($location.path() + '/endpoints');
	};

	$scope.viewUsers = function() {
		$location.path($location.path() + '/users');
	};

	$scope.confirmSave = function(cap) {
		var params = {
			title: 'Update Capability?',
			message: 'Are you sure you want to update the capability?'
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
			save(cap);
		}, function () {
			// do nothing
		});
	};

	$scope.confirmDelete = function(cap) {
		var params = {
			title: 'Delete Capability: ' + cap.name,
			key: cap.name
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
			deleteCapability(cap);
		}, function () {
			// do nothing
		});
	};

};

FormEditCapabilityController.$inject = ['capability', '$scope', '$controller', '$uibModal', '$anchorScroll', '$location', 'locationUtils', 'capabilityService', 'messageModel'];
module.exports = FormEditCapabilityController;
