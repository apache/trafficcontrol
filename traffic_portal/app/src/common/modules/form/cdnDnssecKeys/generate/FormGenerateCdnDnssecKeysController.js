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

/** @typedef {import("moment")} moment */

/**
 * @param {*} cdn
 * @param {*} dnssecKeysRequest
 * @param {*} $scope
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/CDNService")} cdnService
 * @param {import("../../../../models/MessageModel")} messageModel
 */
var FormGenerateCdnDnssecKeysController = function(cdn, dnssecKeysRequest, $scope, $uibModal, formUtils, locationUtils, cdnService, messageModel) {

	var generate = function() {
		$scope.dnssecKeysRequest.effectiveDate = moment($scope.effectiveDate).utc().format();
		cdnService.generateDNSSECKeys($scope.dnssecKeysRequest)
			.then(function(result) {
				messageModel.setMessages(result.data.alerts, true);
				locationUtils.navigateToPath('/cdns/' + cdn.id + '/dnssec-keys');
			});
	};

	$scope.cdn = cdn;
	$scope.dnssecKeysRequest = dnssecKeysRequest;
	$scope.effectiveDate = $scope.dnssecKeysRequest.effectiveDate;

	var ctrl = this;
	ctrl.zeroSeconds = function () {
		if ($scope.effectiveDate) {
			$scope.effectiveDate = $scope.effectiveDate.set({ 'seconds' : 0, });
		}
	};
	$scope.effectiveDate = moment().utc();
	ctrl.zeroSeconds();

	$scope.generateLabel = function() {
		var label = 'Generate DNSSEC Keys';
		if ($scope.ksk_new) {
			label = 'Regenerate DNSSEC Keys';
		}
		return label;
	};

	$scope.msg = 'This will generate DNSSEC keys for the ' + cdn.name + ' CDN and all associated Delivery Services.';

	if ($scope.ksk_new) {
		$scope.msg = 'This will regenerate DNSSEC keys for the ' + cdn.name + ' CDN and all associated Delivery Services. A new DS Record will be created and needs to be added to the parent zone in order for DNSSEC to work properly.';
	}

	$scope.confirmGenerate = function() {
		var title = 'Generate DNSSEC Keys [ ' + cdn.name + ' ]',
			msg = 'This action CANNOT be undone. This will generate DNSSEC keys for the ' + cdn.name + ' CDN and all associated Delivery Services.<br><br>Please type in the name of the CDN to confirm.';

		if ($scope.ksk_new) {
			title = 'Regenerate DNSSEC Keys [ ' + cdn.name + ' ]';
			msg = 'This action CANNOT be undone. This will regenerate DNSSEC keys for the ' + cdn.name + ' CDN and all associated Delivery Services. A new DS Record will be created and needs to be added to the parent zone in order for DNSSEC to work properly.<br><br>Please type in the name of the CDN to confirm.';
		}

		var params = {
			title: title,
			message: msg,
			key: cdn.name
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/enter/dialog.confirm.enter.tpl.html',
			controller: 'DialogConfirmEnterController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			generate();
		}, function () {
			messageModel.setMessages([ { level: 'warning', text: title + ' cancelled' } ], false);
		});
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

};

FormGenerateCdnDnssecKeysController.$inject = ['cdn', 'dnssecKeysRequest', '$scope', '$uibModal', 'formUtils', 'locationUtils', 'cdnService', 'messageModel'];
module.exports = FormGenerateCdnDnssecKeysController;
