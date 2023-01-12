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
 * @param {*} deliveryService
 * @param {*} keys
 * @param {import("../../../../common/api/DeliveryServiceUriSigningKeysService")} deliveryServiceUriSigningKeysService
 * @param {import("../../../../common/models/MessageModel")} messageModel
 * @param {import("../../../../common/service/utils/LocationUtils")} locationUtils
 */
var DeliveryServiceUriSigningKeysController = function($scope, deliveryService, keys, deliveryServiceUriSigningKeysService, messageModel, locationUtils) {
	$scope.deliveryService = deliveryService;
	$scope.keys = keys;
	$scope.keysString = angular.toJson(keys, 4);
	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.saveKeys = function(newKeys) {
		deliveryServiceUriSigningKeysService.setKeys(deliveryService.xmlId, newKeys).then(function() {
			messageModel.setMessages([ { level: 'success', text: 'Keys updated' } ], false);
			$scope.keysString = newKeys;
		}, function() {
			messageModel.setMessages([ { level: 'error', text: 'Failed to update keys, verify that syntax is correct.' } ], false);
		});
	}
};
DeliveryServiceUriSigningKeysController.$inject = ['$scope', 'deliveryService', 'keys', 'deliveryServiceUriSigningKeysService', 'messageModel', 'locationUtils'];
module.exports = DeliveryServiceUriSigningKeysController;
