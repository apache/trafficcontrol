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

var TableCacheGroupParamsUnassignedController = function(cg, allParams, assignedParams, $scope, $uibModalInstance) {

	var selectedParams = [];

	$scope.cg = cg;

	$scope.unassignedParams = allParams.filter(
		function(p) {
			return !assignedParams.has(p.id);
		}
	);

	var addParam = function(paramId) {
		if (selectedParams.indexOf(paramId) === -1) {
			selectedParams.push(paramId);
		}
	};

	var removeParam = function(paramId) {
		selectedParams = selectedParams.filter(
			function (param) {
				return param !== paramId;
			}
		);
	};

	$scope.updateParams = function($event, paramId) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addParam(paramId);
		} else {
			removeParam(paramId);
		}
	};

	$scope.handleRowClick = function($index) {
		$('#checkbox-' + $index).trigger('click');
	};

	$scope.submit = function() {
		$uibModalInstance.close(selectedParams);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	angular.element(document).ready(function () {
		$('#cgParamsUnassignedTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"order": [[ 1, 'asc' ]],
			"columnDefs": [
				{ "width": "5%", "targets": 0 },
				{ "width": "50%", "targets": 3 }
			]
		});
	});

};

TableCacheGroupParamsUnassignedController.$inject = ['cg', 'allParams', 'assignedParams', '$scope', '$uibModalInstance'];
module.exports = TableCacheGroupParamsUnassignedController;
