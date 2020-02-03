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

var TableCacheGroupParametersController = function(cacheGroup, parameters, $controller, $scope, $state, $uibModal, locationUtils, cacheGroupParameterService) {

	// extends the TableParametersController to inherit common methods
	angular.extend(this, $controller('TableParametersController', { parameters: parameters, $scope: $scope }));

	$scope.cacheGroup = cacheGroup;

	// adds some items to the base parameters context menu
	$scope.contextMenuItems.splice(2, 0,
		{
			text: 'Unlink Parameter from Cache Group',
			click: function ($itemScope) {
				$scope.removeParameter($itemScope.p.id);
			}
		},
		null // Divider
	);

	$scope.removeParameter = function(paramId) {
		cacheGroupParameterService.unlinkCacheGroupParameter(cacheGroup.id, paramId)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.selectParams = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/cacheGroupParameters/table.cacheGroupParamsUnassigned.tpl.html',
			controller: 'TableCacheGroupParamsUnassignedController',
			size: 'lg',
			resolve: {
				cg: function() {
					return cacheGroup;
				},
				parameters: function(parameterService) {
					return parameterService.getCacheGroupUnassignedParams(cacheGroup.id);
				}
			}
		});
		modalInstance.result.then(function(selectedParams) {
			var massagedArray = [];
			for (i = 0; i < selectedParams.length; i++) {
				massagedArray.push( { cacheGroupId: cacheGroup.id, parameterId: selectedParams[i] } );
			}
			cacheGroupParameterService.linkCacheGroupParameters(massagedArray)
				.then(
					function() {
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#cacheGroupParametersTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ 'orderable': false, 'targets': 4 },
				{ "width": "50%", "targets": 2 }
			],
			"aaSorting": []
		});
	});

};

TableCacheGroupParametersController.$inject = ['cacheGroup', 'parameters', '$controller', '$scope', '$state', '$uibModal', 'locationUtils', 'cacheGroupParameterService'];
module.exports = TableCacheGroupParametersController;
