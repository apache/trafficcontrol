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

var FormDeliveryServiceSslKeysController = function(deliveryService, sslKeys, $scope, locationUtils, deliveryServiceSslKeysService, $uibModal, $anchorScroll, formUtils, $filter) {

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

	$scope.deliveryService = deliveryService;
	$scope.sslKeys = setSSLKeys(sslKeys);
	if ($scope.sslKeys.authType === undefined || $scope.sslKeys.authType === '') {
        $scope.sslKeys.authType = 'Not Assigned';
    }

	$scope.hasError = formUtils.hasError;
	$scope.hasPropertyError = formUtils.hasPropertyError;
	$scope.navigateToPath = locationUtils.navigateToPath;

	$scope.formattedExpiration = $scope.sslKeys.expiration !== undefined ? $filter('date')($scope.sslKeys.expiration, 'MM/dd/yyyy') : undefined;

	$scope.generateKeys = function() {
		locationUtils.navigateToPath('/delivery-services/' + deliveryService.id + '/ssl-keys/generate');
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
                    if ($scope.dsSslKeyForm) $scope.dsSslKeyForm.$setPristine();
                });
		});
	};

};

FormDeliveryServiceSslKeysController.$inject = ['deliveryService', 'sslKeys', '$scope', 'locationUtils', 'deliveryServiceSslKeysService', '$uibModal', '$anchorScroll', 'formUtils', '$filter'];
module.exports = FormDeliveryServiceSslKeysController;
