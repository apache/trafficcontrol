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

var FormCdnDnssecKeysController = function(cdn, dnssecKeys, $scope, $location, $uibModal, dateUtils, formUtils, stringUtils, locationUtils, messageModel) {

	var generate = function() {
		$location.path($location.path() + '/generate');
	};
	var regenerateKSK = function() {
		$location.path($location.path() + '/regenerateKsk');
	};

	$scope.cdn = cdn;

	$scope.ksk_new = null;

	$scope.falseTrue = [
		{ value: true, label: 'true' },
		{ value: false, label: 'false' }
	];

	$scope.generateLabel = function() {
		var label = 'Generate DNSSEC Keys';
		if ($scope.ksk_new) {
			label = 'Regenerate DNSSEC Keys';
		}
		return label;
	};

	$scope.confirmGenerate = function() {
		var title = 'Generate DNSSEC Keys [ ' + cdn.name + ' ]',
			msg = 'This will generate DNSSEC keys for the ' + cdn.name + ' CDN and all associated Delivery Services.<br><br>Are you sure you want to proceed?';

		if ($scope.ksk_new) {
			title = 'Regenerate DNSSEC Keys [ ' + cdn.name + ' ]';
			msg = 'This will regenerate DNSSEC keys for the ' + cdn.name + ' CDN and all associated Delivery Services. A new DS Record will be created and needs to be added to the parent zone in order for DNSSEC to work properly.<br><br>Are you sure you want to proceed?';
		}

		var params = {
			title: title,
			message: msg
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
			generate();
		}, function () {
			messageModel.setMessages([ { level: 'warning', text: title + ' cancelled' } ], false);
		});
	};

	$scope.confirmKSK = function() {
		var title = 'Regenerate KSK Keys [ ' + cdn.name + ' ]',
			msg = 'This will regenerate KSK keys for the ' + cdn.name + ' CDN and all associated Delivery Services.<br><br>Are you sure you want to proceed?';

		var params = {
			title: title,
			message: msg
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
			regenerateKSK();
		}, function () {
			messageModel.setMessages([ { level: 'warning', text: title + ' cancelled' } ], false);
		});
	};


	$scope.dateFormat = dateUtils.dateFormat;

	$scope.navigateToPath = locationUtils.navigateToPath;

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

	var init = function() {
		if (dnssecKeys[cdn.name]) {
			$scope.ksk_new = _.find(dnssecKeys[cdn.name].ksk, function(ksk) { return ksk.status == 'new' });
		}
	};
	init();

};

FormCdnDnssecKeysController.$inject = ['cdn', 'dnssecKeys', '$scope', '$location', '$uibModal', 'dateUtils', 'formUtils', 'stringUtils', 'locationUtils', 'messageModel'];
module.exports = FormCdnDnssecKeysController;
