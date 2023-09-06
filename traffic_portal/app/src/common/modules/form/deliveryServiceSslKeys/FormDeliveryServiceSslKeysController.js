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
 * @param {*} deliveryService
 * @param {*} sslKeys
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/DeliveryServiceSslKeysService")} deliveryServiceSslKeysService
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("angular").IFilterService} $filter
 */
var FormDeliveryServiceSslKeysController = function(deliveryService, sslKeys, $scope, $state, locationUtils, deliveryServiceSslKeysService, $uibModal, $anchorScroll, formUtils, $filter) {

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.showRevert = function() {
		return $scope.sslKeys.version > 1
	};

	var setSSLKeys = function(sslKeys) {
		if (!sslKeys.hostname) {
			var url = deliveryService.exampleURLs[0],
				defaultHostName = url.split("://")[1];
			if (deliveryService.type.indexOf('HTTP') != -1) {
				var parts = defaultHostName.split(".");
				parts[0] = "*";
				defaultHostName = parts.join(".");
			}
			sslKeys.hostname = defaultHostName;
		}
		return sslKeys;
	};

	var getAcmeProviders = function() {
		deliveryServiceSslKeysService.getAcmeProviders()
			.then(function(result) {
				$scope.acmeProviders = result;
				if (!$scope.acmeProviders.includes('Lets Encrypt')) {
					$scope.acmeProviders.push('Lets Encrypt');
				}
				if (!$scope.acmeProviders.includes('Self Signed')) {
					$scope.acmeProviders.push('Self Signed');
				}
				if (!$scope.acmeProviders.includes('Not Assigned')) {
					$scope.acmeProviders.push('Not Assigned');
				}
				if (!$scope.acmeProviders.includes('Provided Manually')) {
					$scope.acmeProviders.push('Provided Manually');
				}
			});
	};

	$scope.acmeProviders = [];
	$scope.deliveryService = deliveryService;
	$scope.sslKeys = setSSLKeys(sslKeys);
	if ($scope.sslKeys.authType === undefined || $scope.sslKeys.authType === '') {
        $scope.sslKeys.authType = 'Not Assigned';
    }
	$scope.acmeProvider = sslKeys.authType;

	$scope.hasError = formUtils.hasError;
	$scope.hasPropertyError = formUtils.hasPropertyError;
	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.formattedExpiration = $scope.sslKeys.expiration !== undefined ? $filter('date')($scope.sslKeys.expiration, 'MM/dd/yyyy') : undefined;
	$scope.sans = $scope.sslKeys.sans !== undefined ? sslKeys.sans.join(', ') : ""

	$scope.generateKeys = function() {
		locationUtils.navigateToPath('/delivery-services/' + deliveryService.id + '/ssl-keys/generate');
	};

	$scope.navState = 0;
	$scope.updateState = function(newState) {
		$scope.navState = newState;
	}

	$scope.renewCert = function() {
		var params = {
			title: 'Renew SSL Keys for Delivery Service: ' + deliveryService.xmlId
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
			deliveryServiceSslKeysService.renewCert(deliveryService).then(
				function() {
					$anchorScroll();
					$scope.refresh();
					if ($scope.dsSslKeyForm) $scope.dsSslKeyForm.$setPristine();
				});
		});
	};

	$scope.save = function() {
		var params = {
			title: 'Add New SSL Keys for Delivery Service: ' + deliveryService.xmlId
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
			deliveryServiceSslKeysService.addSslKeys(sslKeys, deliveryService).then(
                function() {
                    $anchorScroll();
                    $scope.refresh();
                    if ($scope.dsSslKeyForm) $scope.dsSslKeyForm.$setPristine();
                });
		});
	};

	$scope.confirmDelete = function() {
		var params = {
			title: 'Delete latest SSL Keys for Delivery Service: ' + deliveryService.xmlId,
			key: deliveryService.xmlId + '-keys'
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
			deliveryServiceSslKeysService.deleteCert(deliveryService).then(
				function() {
					$anchorScroll();
					$scope.refresh();
				});
		});
	};

	$scope.updateProvider = function() {
		sslKeys.authType = $scope.acmeProvider;
	};

	var init = function () {
		getAcmeProviders();
	};
	init();

	$scope.revertCert = function() {
		var params = {
			title: 'Revert SSL Keys for Delivery Service: ' + deliveryService.xmlId,
			message: 'This will replace existing keys with the keys from the previous version'

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
			deliveryServiceSslKeysService.revertCert(deliveryService).then(
				function() {
					$anchorScroll();
					$scope.refresh();
					if ($scope.dsSslKeyForm) $scope.dsSslKeyForm.$setPristine();
				});
		});
	};

};

FormDeliveryServiceSslKeysController.$inject = ['deliveryService', 'sslKeys', '$scope', '$state', 'locationUtils', 'deliveryServiceSslKeysService', '$uibModal', '$anchorScroll', 'formUtils', '$filter'];
module.exports = FormDeliveryServiceSslKeysController;
