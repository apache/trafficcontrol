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

var TableProfileParamsUnassignedController = function(profile, allParams, assignedParams, $scope, $uibModalInstance) {

	var selectedParamIds = [];

	var addParam = function(paramId) {
		if (_.indexOf(selectedParamIds, paramId) == -1) {
			selectedParamIds.push(paramId);
		}
	};

	var removeParam = function(paramId) {
		selectedParamIds = _.without(selectedParamIds, paramId);
	};

	var addAll = function() {
		markParams(true);
		selectedParamIds = _.pluck(allParams, 'id');
	};

	var removeAll = function() {
		markParams(false);
		selectedParamIds = [];
	};

	var markParams = function(selected) {
		$scope.selectedParams = _.map(allParams, function(param) {
			param['selected'] = selected;
			return param;
		});
	};

	$scope.profile = profile;

	$scope.selectedParams = _.map(allParams, function(param) {
		var isAssigned = _.find(assignedParams, function(assignedParam) { return assignedParam.id == param.id });
		if (isAssigned) {
			param['selected'] = true; // so the checkbox will be checked
			selectedParamIds.push(param.id); // so the param is added to selected params
		}
		return param;
	});

	$scope.allSelected = function() {
		return allParams.length == selectedParamIds.length;
	};

	$scope.selectAll = function($event) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addAll();
		} else {
			removeAll();
		}
	};

	$scope.updateParams = function($event, paramId) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addParam(paramId);
		} else {
			removeParam(paramId);
		}
	};

	$scope.submit = function() {
		$uibModalInstance.close(selectedParamIds);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	angular.element(document).ready(function () {
		var profileParamsUnassignedTable = $('#profileParamsUnassignedTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"order": [[ 1, 'asc' ]],
			"columnDefs": [
				{ "width": "5%", "targets": 0 }
			],
			"stateSave": false
		});
		profileParamsUnassignedTable.on( 'search.dt', function () {
			var search = $('#profileParamsUnassignedTable_filter input').val();
			if (search.length > 0) {
				$("#selectAllCB").attr("disabled", true);
			} else {
				$("#selectAllCB").removeAttr("disabled");
			}
		} );
	});

};

TableProfileParamsUnassignedController.$inject = ['profile', 'allParams', 'assignedParams', '$scope', '$uibModalInstance'];
module.exports = TableProfileParamsUnassignedController;
