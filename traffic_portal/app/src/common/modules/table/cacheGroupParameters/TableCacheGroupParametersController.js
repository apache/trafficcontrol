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

	let cacheGroupParametersTable;

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
				allParams: function(parameterService) {
					return parameterService.getParameters();
				},
				assignedParams: function() {
					return new Set(parameters.map(function(x){return x.id;}));
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

	$scope.toggleVisibility = function(colName) {
		const col = cacheGroupParametersTable.column(colName + ':name');
		col.visible(!col.visible());
		cacheGroupParametersTable.rows().invalidate().draw();
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		cacheGroupParametersTable = $('#cacheGroupParametersTable').DataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": [],
			"columnDefs": [
				{ "width": "50%", "targets": 2 },
				{ "orderable": false, "targets": 4 }
			],
			"columns": [
				{ "name": "Name", "visible": true, "searchable": true },
				{ "name": "Config File", "visible": true, "searchable": true },
				{ "name": "Value", "visible": true, "searchable": true },
				{ "name": "Secure", "visible": true, "searchable": true },
				{ "name": "Action", "visible": true, "searchable": false }
			],
			"initComplete": function(settings, json) {
				try {
					// need to create the show/hide column checkboxes and bind to the current visibility
					$scope.columns = JSON.parse(localStorage.getItem('DataTables_cacheGroupParametersTable_/')).columns;
				} catch (e) {
					console.error("Failure to retrieve required column info from localStorage (key=DataTables_cacheGroupParametersTable_/):", e);
				}
			}
		});
	});

};

TableCacheGroupParametersController.$inject = ['cacheGroup', 'parameters', '$controller', '$scope', '$state', '$uibModal', 'locationUtils', 'cacheGroupParameterService'];
module.exports = TableCacheGroupParametersController;
