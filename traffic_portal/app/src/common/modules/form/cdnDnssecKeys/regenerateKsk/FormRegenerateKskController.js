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

var FormRegenerateKskController = function(cdn, dnssecKeysRequest, $scope, $location, $uibModal, formUtils, locationUtils, cdnService, messageModel) {

	var generate = function() {
		$scope.kskRequest.effectiveDate = moment($scope.kskRequest.effectiveDate).utc().format();
		cdnService.regenerateKSK($scope.kskRequest, $scope.cdnKey)
			.then(function(result) {
				messageModel.setMessages(result.data.alerts, true);
				locationUtils.navigateToPath('/cdns/' + cdn.id + '/dnssec-keys');
			});
	};

	$scope.cdn = cdn;
	$scope.cdnKey = dnssecKeysRequest.key;
	$scope.kskRequest = {
		effectiveDate: dnssecKeysRequest.effectiveDate,
		expirationDays: 365
	};

	$scope.generateLabel = function() {
		var label = 'Regenerate KSK';
		return label;
	};

	$scope.confirmGenerate = function() {
		var title = 'Regenerate DNSSEC Keys [ ' + cdn.name + ' ]',
			msg = 'This action CANNOT be undone. This will regenerate KSK keys for the ' + cdn.name + ' CDN and all associated Delivery Services.<br><br>Please type in the name of the CDN to confirm.';

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

	$scope.navigateToPath = locationUtils.navigateToPath;

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

};

FormRegenerateKskController.$inject = ['cdn', 'dnssecKeysRequest', '$scope', '$location', '$uibModal', 'formUtils', 'locationUtils', 'cdnService', 'messageModel'];
module.exports = FormRegenerateKskController;
