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
 * @param {*} regex
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/DeliveryServiceRegexService")} deliveryServiceRegexService
 */
var FormEditDeliveryServiceRegexController = function(deliveryService, regex, $scope, $controller, $uibModal, $anchorScroll, locationUtils, deliveryServiceRegexService) {

	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceRegexController', { deliveryService: deliveryService, regex: regex, $scope: $scope }));

	var deleteDeliveryServiceRegex = function(dsId, regexId) {
		deliveryServiceRegexService.deleteDeliveryServiceRegex(dsId, regexId)
			.then(function() {
				locationUtils.navigateToPath('/delivery-services/' + dsId + '/regexes');
			});
	};

	$scope.regexPattern = angular.copy(regex.pattern);

	$scope.settings = {
		isNew: false,
		saveLabel: 'Update'
	};

	$scope.save = function(dsId, regex) {
		deliveryServiceRegexService.updateDeliveryServiceRegex(dsId, regex).
			then(function() {
				$scope.regexPattern = angular.copy(regex.pattern);
				$anchorScroll(); // scrolls window to top
			});
	};

	$scope.confirmDelete = function(regex) {
		var params = {
			title: 'Delete Delivery Service Regex: ' + regex.pattern,
			key: regex.pattern
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
			deleteDeliveryServiceRegex(deliveryService.id, regex.id);
		}, function () {
			// do nothing
		});
	};

};

FormEditDeliveryServiceRegexController.$inject = ['deliveryService', 'regex', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'deliveryServiceRegexService'];
module.exports = FormEditDeliveryServiceRegexController;
