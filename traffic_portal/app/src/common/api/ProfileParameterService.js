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

var ProfileParameterService = function(Restangular, httpService, messageModel, ENV, $uibModal) {

	this.unlinkProfileParameter = function(profileId, paramId) {
		return httpService.delete(ENV.api['root'] + 'profileparameters/' + profileId + '/' + paramId)
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Profile and parameter were unlinked.' } ], false);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, true);
				}
			);
	};

	this.linkProfileParameters = function(profileId, params) {
		return Restangular.service('profileparameter').post({ profileId: profileId, paramIds: params, replace: true })
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Parameters linked to profile' } ], false);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

	var linkParamProfilesHelper = function(paramId, profiles) { 
		return Restangular.service('parameterprofile').post({ paramId: paramId, profileIds: profiles, replace: true })
	.then(
		function() {
			messageModel.setMessages([ { level: 'success', text: 'Profiles linked to parameter' } ], false);
		},
		function(fault) {
			messageModel.setMessages(fault.data.alerts, false);
		}
	);
	}

	this.linkParamProfiles = function(paramId, profiles) {
		return linkParamProfilesHelper(paramId, profiles);
	};

	this.selectProfiles = function(parameter, profiles) {
		return $uibModal.open({
		templateUrl: 'common/modules/table/parameterProfiles/table.paramProfilesUnassigned.tpl.html',
		controller: 'TableParamProfilesUnassignedController',
		size: 'lg',
		resolve: {
			parameter: function() {
				return parameter;
			},
			allProfiles: function(profileService) {
				return profileService.getProfiles({ orderby: 'name' });
			},
			assignedProfiles: function(profileService) {
				return profiles || profileService.getParameterProfiles(parameter.id); // there's an uncaught error that doesn't affect functionality if the parameter creation fails
			}
		}
		}).result.then(function(selectedProfileIds) {
		var params = {
			title: 'Assign profiles to ' + parameter.name,
			message: 'Are you sure you want to modify the profiles assigned to ' + parameter.name + '?'
		};
		return $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		}).result.then(function() {
			linkParamProfilesHelper(parameter.id, selectedProfileIds); // not ideal, but it's what works for now
		});
	});
};

};

ProfileParameterService.$inject = ['Restangular', 'httpService', 'messageModel', 'ENV', '$uibModal'];
module.exports = ProfileParameterService;
