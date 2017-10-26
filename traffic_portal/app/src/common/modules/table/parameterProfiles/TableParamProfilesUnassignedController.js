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

var TableParamProfilesUnassignedController = function(parameter, allProfiles, assignedProfiles, $scope, $uibModalInstance) {

	var selectedProfileIds = [];

	var addProfile = function(profileId) {
		if (_.indexOf(selectedProfileIds, profileId) == -1) {
			selectedProfileIds.push(profileId);
		}
	};

	var removeProfile = function(profileId) {
		selectedProfileIds = _.without(selectedProfileIds, profileId);
	};

	var addAll = function() {
		markProfiles(true);
		selectedProfileIds = _.pluck(allProfiles, 'id');
	};

	var removeAll = function() {
		markProfiles(false);
		selectedProfileIds = [];
	};

	var markProfiles = function(selected) {
		$scope.selectedProfiles = _.map(allProfiles, function(profile) {
			profile['selected'] = selected;
			return profile;
		});
	};

	$scope.parameter = parameter;

	$scope.selectedProfiles = _.map(allProfiles, function(profile) {
		var isAssigned = _.find(assignedProfiles, function(assignedProfile) { return assignedProfile.id == profile.id });
		if (isAssigned) {
			profile['selected'] = true; // so the checkbox will be checked
			selectedProfileIds.push(profile.id); // so the profile is added to selected profiles
		}
		return profile;
	});

	$scope.allSelected = function() {
		return allProfiles.length == selectedProfileIds.length;
	};

	$scope.selectAll = function($event) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addAll();
		} else {
			removeAll();
		}
	};

	$scope.updateProfiles = function($event, profileId) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addProfile(profileId);
		} else {
			removeProfile(profileId);
		}
	};

	$scope.submit = function() {
		$uibModalInstance.close(selectedProfileIds);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	angular.element(document).ready(function () {
		var paramProfilesUnassignedTable = $('#paramProfilesUnassignedTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"order": [[ 1, 'asc' ]],
			"columnDefs": [
				{ "width": "5%", "targets": 0 }
			],
			"stateSave": false
		});
		paramProfilesUnassignedTable.on( 'search.dt', function () {
			var search = $('#paramProfilesUnassignedTable_filter input').val();
			if (search.length > 0) {
				$("#selectAllCB").attr("disabled", true);
			} else {
				$("#selectAllCB").removeAttr("disabled");
			}
		} );
	});

};

TableParamProfilesUnassignedController.$inject = ['parameter', 'allProfiles', 'assignedProfiles', '$scope', '$uibModalInstance'];
module.exports = TableParamProfilesUnassignedController;
