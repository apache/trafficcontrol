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

var TableProfilesParamsCompareController = function(profile1, profile2, profilesParams, $scope, $state, $q, $uibModal, messageModel, locationUtils, deliveryServiceService, profileParameterService, serverService) {

	this.profile1Usage;
	this.profile2Usage;

	let updateProfile1 = false,
		updateProfile2 = false;

	let getProfileUsage = function(profile, profNum) {
		if (profile.type === 'DS_PROFILE') { // if this is a ds profile, then it is used by delivery service(s) so we'll fetch the ds count...
			deliveryServiceService.getDeliveryServices({ profile: profile.id }).
				then(function(result) {
					this['profile' + profNum + 'Usage'] = result.length + ' delivery services';
				});
		} else { // otherwise the profile is used by servers so we'll fetch the server count...
			serverService.getServers({ profileId: profile.id }).
				then(function(result) {
					this['profile' + profNum + 'Usage'] = result.length + ' servers';
				});
		}
	};

	let confirmUpdateProfiles = function() {
		// ok, this method is fun :)
		let params = {
			title: 'Modify ' + profile1.name + ' parameters',
			message: 'The ' + profile1.name + ' profile is used by ' + this.profile1Usage + '. Are you sure you want to modify the parameters?'
		};
		let modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result
			.then(
				function() { updateProfile1 = true; },
				function() { updateProfile1 = false; }
			)
			.finally(
				function() {
					let params = {
						title: 'Modify ' + profile2.name + ' parameters',
						message: 'The ' + profile2.name + ' profile is used by ' + this.profile2Usage + '. Are you sure you want to modify the parameters?'
					};
					let modalInstance = $uibModal.open({
						templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
						controller: 'DialogConfirmController',
						size: 'md',
						resolve: {
							params: function () {
								return params;
							}
						}
					});
					modalInstance.result
						.then(
							function() { updateProfile2 = true; },
							function() { updateProfile2 = false; }
						)
						.finally(
							function() {
								let promises = [];

								if (updateProfile1) {
									promises.push(profileParameterService.linkProfileParameters(profile1, _.pluck($scope.selectedProfile1Params, 'id')));
								}

								if (updateProfile2) {
									promises.push(profileParameterService.linkProfileParameters(profile2, _.pluck($scope.selectedProfile2Params, 'id')));
								}

								if (promises.length > 0) {
									$q.all(promises)
										.then(
											function(result) {
												let messages = [];
												for (let i = 0; i < result.length; i++) {
													messages = _.union(messages, result[i].data.alerts);
												};
												messageModel.setMessages(messages, false);
											})
										.finally(
											function() {
												$scope.refresh();
											}
										);
								}
							}
						);
				}
			);
	};

	$scope.dirty = false;

	$scope.profile1 = profile1;
	$scope.profile2 = profile2;
	$scope.profilesParams = profilesParams;

	$scope.selectedProfile1Params = [];
	$scope.selectedProfile2Params = [];

	$scope.selectedParams = _.map(profilesParams, function(param) {
		let isAssignedToProfile1 = _.find(profile1.params, function(profile1param) { return profile1param.id == param.id }),
			isAssignedToProfile2 = _.find(profile2.params, function(profile2param) { return profile2param.id == param.id });

		if (isAssignedToProfile1) {
			param['origSelected1'] = true;
			param['selected1'] = true;
		} else {
			param['origSelected1'] = false;
			param['selected1'] = false;
		}

		if (isAssignedToProfile2) {
			param['origSelected2'] = true;
			param['selected2'] = true;
		} else {
			param['origSelected2'] = false;
			param['selected2'] = false;
		}

		return param;
	});

	$scope.updateSelectedCount = function(profNum) {
		$scope['selectedProfile' + profNum + 'Params'] = _.filter($scope.selectedParams, function(param) { return param['selected' + profNum] == true; } );
	};

	$scope.isCheckboxDirty = function(pp, profNum) {
		return pp['selected' + profNum] !== pp['origSelected' + profNum];
	};

	$scope.onChange = function(profNum) {
		$scope.dirty = true;
		$scope.updateSelectedCount(profNum);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.restore = function() {
		let params = {
			title: 'Restore Parameter Assignments?',
			message: 'Any changes you have made will be lost and the original parameter assignments will be restored.'
		};
		let modalInstance = $uibModal.open({
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
			$scope.refresh();
		});
	};

	$scope.update = function() {
		confirmUpdateProfiles();
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#profilesParamsCompareTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": -1,
			"aaSorting": [],
			"columnDefs": [
				{ "width": "50%", "targets": 2 }
			]
		});

		getProfileUsage(profile1, 1);
		getProfileUsage(profile2, 2);

		$scope.updateSelectedCount(1);
		$scope.updateSelectedCount(2);
	});

};

TableProfilesParamsCompareController.$inject = ['profile1', 'profile2', 'profilesParams', '$scope', '$state', '$q', '$uibModal', 'messageModel', 'locationUtils', 'deliveryServiceService', 'profileParameterService', 'serverService'];
module.exports = TableProfilesParamsCompareController;
