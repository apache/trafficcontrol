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

var ProfileParameterService = function($http, messageModel, ENV) {

	this.unlinkProfileParameter = function(profileId, paramId) {
		return $http.delete(ENV.api.unstable + 'profileparameters/' + profileId + '/' + paramId).then(
				function(result) {
					messageModel.setMessages([ { level: 'success', text: 'Profile and parameter were unlinked.' } ], false);
					return result;
				},
				function(err) {
					messageModel.setMessages(err.data.alerts, true);
					throw err;
				}
			);
	};

	this.linkProfileParameters = function(profile, params) {
		return $http.post(ENV.api.unstable + 'profileparameter', { profileId: profile.id, paramIds: params, replace: true }).then(
			function(result) {
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.linkParamProfiles = function(paramId, profiles) {
		return $http.post(ENV.api.unstable + 'parameterprofile', { paramId: paramId, profileIds: profiles, replace: true }).then(
			function(result) {
				messageModel.setMessages([ { level: 'success', text: 'Profiles linked to parameter' } ], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

};

ProfileParameterService.$inject = ['$http', 'messageModel', 'ENV'];
module.exports = ProfileParameterService;
