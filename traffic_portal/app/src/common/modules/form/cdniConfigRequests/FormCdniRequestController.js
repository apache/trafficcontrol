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
 * @param {*} $scope
 * @param {*} $stateParams
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../api/CdniService")} cdniService
 * @param {*} cdniRequest
 * @param {*} currentConfig
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../models/MessageModel")} messageModel
 */
var FormCdniRequestController = function($scope, $stateParams, $uibModal, cdniService, cdniRequest, currentConfig, locationUtils, messageModel) {
	$scope.reqId = $stateParams.reqId;
	$scope.cdniRequest = cdniRequest;
	$scope.cdniRequest.data = JSON.stringify($scope.cdniRequest.data, null, 5);
	$scope.currentConfig = JSON.stringify(currentConfig, null, 5);

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.respondToRequest = function(approve) {
		const titleStart = approve ? 'Approve' : 'Deny';
		const params = {
			title: `${titleStart} CDNi Update Request: ${cdniRequest.id}`
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {params}
		});
		modalInstance.result.then(function() {
			cdniService.sendResponseToCdniRequest(cdniRequest.id, approve).then(
				function(result) {
					messageModel.setMessages([{level: 'success', text: result}], true);
					$scope.navigateToPath('/cdni-config-requests')
				});
		});
	};
};

FormCdniRequestController.$inject = ['$scope', '$stateParams', '$uibModal', 'cdniService', 'cdniRequest', 'currentConfig', 'locationUtils', 'messageModel'];
module.exports = FormCdniRequestController;
