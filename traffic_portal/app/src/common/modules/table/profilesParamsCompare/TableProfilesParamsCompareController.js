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

/** @typedef {import("jquery")} */

/**
 * This is a controller for the table used to compare the Parameters of two
 * Profiles.
 *
 * @param {{name: string; params: {id: number}[]}} profile1
 * @param {{name: string; params: {id: number}[]}} profile2
 * @param {{id: number}[]} profilesParams
 * @param {boolean} showAll
 * @param {*} $scope
 * @param {*} $state
 * @param {import("angular").IQService} $q
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../models/MessageModel")} messageModel
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/DeliveryServiceService")} deliveryServiceService
 * @param {import("../../../api/ProfileParameterService")} profileParameterService
 * @param {import("../../../api/ServerService")} serverService
 */
var TableProfilesParamsCompareController = function(profile1, profile2, profilesParams, showAll, $scope, $state, $q, $uibModal, messageModel, locationUtils, deliveryServiceService, profileParameterService, serverService) {

	let updateProfile1 = false,
		updateProfile2 = false;

	async function getProfileUsage(profile, profNum) {
		if (profile.type === 'DS_PROFILE') { // if this is a ds profile, then it is used by delivery service(s) so we'll fetch the ds count...
			const result = await deliveryServiceService.getDeliveryServices({ profile: profile.id });
			$scope[`profile${profNum}Usage`] = `${result.length} delivery services`;
		} else { // otherwise the profile is used by servers so we'll fetch the server count...
			const result = await serverService.getServers({ profileName: profile.name });
			$scope[`profile${profNum}Usage`] = `${result.length} servers`;
		}
	};

	$scope.showAll = showAll;

	$scope.dirty = false;

	$scope.profile1Usage;
	$scope.profile2Usage;

	$scope.profile1 = profile1;
	$scope.profile2 = profile2;
	$scope.profilesParams = profilesParams;

	$scope.selectedProfile1Params = [];
	$scope.selectedProfile2Params = [];

	$scope.confirmUpdateProfiles = function() {
		// ok, this method is fun :)
		let params = {
			title: 'Modify ' + profile1.name + ' parameters',
			message: 'The ' + profile1.name + ' profile is used by ' + $scope.profile1Usage + '. Are you sure you want to modify the parameters?'
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
						message: 'The ' + profile2.name + ' profile is used by ' + $scope.profile2Usage + '. Are you sure you want to modify the parameters?'
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
									const selectedProfile1ParamIds = $scope.selectedProfile1Params.map(function(pp){return pp.id;});
									promises.push(profileParameterService.linkProfileParameters(profile1, selectedProfile1ParamIds));
								}

								if (updateProfile2) {
									const selectedProfile2ParamIds = $scope.selectedProfile2Params.map(function(pp){return pp.id;});
									promises.push(profileParameterService.linkProfileParameters(profile2, selectedProfile2ParamIds));
								}

								if (promises.length > 0) {
									$q.all(promises)
										.then(
											function(result) {
												const messages = new Set();
												for (let i = 0; i < result.length; i++) {
													for (let j = 0; j < result[i].data.alerts.length; j++) {
														messages.add(result[i].data.alerts[j]);
													};
												};
												messageModel.setMessages(Array.from(messages), false);
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

	$scope.selectedParams = profilesParams.map(function(param) {
		let isAssignedToProfile1 = profile1.params ? profile1.params.some(pp => pp.id === param.id) : false,
			isAssignedToProfile2 = profile2.params ? profile2.params.some(pp => pp.id === param.id) : false;

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
		const num = parseInt(profNum, 10);
		$scope['selectedProfile' + num + 'Params'] = $scope.selectedParams.filter(function(p){return p['selected' + num] == true;});
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

	$scope.filterFn = function(pp) {
		if ($scope.showAll) return true;
		return (pp.selected1 !== pp.selected2);
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	angular.element(document).ready(function () {
		// Datatables should be replaced with AG-Grid.
		// @ts-ignore
		$('#profilesParamsCompareTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": -1,
			"aaSorting": [],
			"columnDefs": [
				{ "width": "50%", "targets": 2 }
			],
			"language": {
				"emptyTable": ($scope.showAll) ? "No data available in table" : "Profiles are identical"
			}
		});

		getProfileUsage(profile1, 1);
		getProfileUsage(profile2, 2);

		$scope.updateSelectedCount(1);
		$scope.updateSelectedCount(2);
	});

};

TableProfilesParamsCompareController.$inject = ['profile1', 'profile2', 'profilesParams', 'showAll', '$scope', '$state', '$q', '$uibModal', 'messageModel', 'locationUtils', 'deliveryServiceService', 'profileParameterService', 'serverService'];
module.exports = TableProfilesParamsCompareController;
