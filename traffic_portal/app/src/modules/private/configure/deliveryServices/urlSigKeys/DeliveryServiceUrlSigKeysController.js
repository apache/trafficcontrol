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

var DeliveryServiceUrlSigKeysController = function(deliveryService, urlSigKeys, $scope, $state, locationUtils, deliveryServiceService, deliveryServiceUrlSigKeysService, $uibModal) {
	$scope.deliveryService = deliveryService;
	$scope.urlSigKeys = Object.keys(urlSigKeys).map(function(key) {
			return {sortBy: parseInt(key.slice(3)), label: key, urlSigKey: urlSigKeys[key]};
	});

	$scope.generateUrlSigKeys = function() {
		deliveryServiceUrlSigKeysService.generateUrlSigKeys(deliveryService.xmlId).then(
			function() {
			$state.reload();
		});
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.selectCopyFromDS = function() {
        var params = {
            title: 'Copy URL Sig Keys to: ' + deliveryService.displayName,
            message: "Please select a Delivery Service to copy from:",
            label: "xmlId"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
            controller: 'DialogSelectController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function(deliveryServiceService) {
                    return deliveryServiceService.getDeliveryServices({ signed: true })
                    .then(function(result){
                    	return _.filter(result, function(ds){
                    		return ds.id !== deliveryService.id;
                    	})
                    });
                }
            }
        });
        modalInstance.result.then(function(copyFromDs) {
            deliveryServiceUrlSigKeysService.copyUrlSigKeys(deliveryService.xmlId, copyFromDs.xmlId)
            .then(
           		function() {
            	$state.reload();
        	});
        });
    };

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#regexesTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});
};

DeliveryServiceUrlSigKeysController.$inject = ['deliveryService', 'urlSigKeys', '$scope', '$state', 'locationUtils', 'deliveryServiceService', 'deliveryServiceUrlSigKeysService', '$uibModal'];
module.exports = DeliveryServiceUrlSigKeysController;
