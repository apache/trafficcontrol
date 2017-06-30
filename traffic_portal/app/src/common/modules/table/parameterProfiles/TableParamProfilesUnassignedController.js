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

var TableParamProfilesUnassignedController = function(parameter, profiles, $scope, $uibModalInstance) {

	var selectedProfiles = [];

	$scope.parameter = parameter;

	$scope.unassignedProfiles = profiles;

	var addProfile = function(profileId) {
		if (_.indexOf(selectedProfiles, profileId) == -1) {
			selectedProfiles.push(profileId);
		}
	};

	var removeProfile = function(profileId) {
		selectedProfiles = _.without(selectedProfiles, profileId);
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
		$uibModalInstance.close(selectedProfiles);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	angular.element(document).ready(function () {
		$('#paramProfilesUnassignedTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"order": [[ 1, 'asc' ]],
			"columnDefs": [
				{ "width": "5%", "targets": 0 }
			]
		});
	});

};

TableParamProfilesUnassignedController.$inject = ['parameter', 'profiles', '$scope', '$uibModalInstance'];
module.exports = TableParamProfilesUnassignedController;
