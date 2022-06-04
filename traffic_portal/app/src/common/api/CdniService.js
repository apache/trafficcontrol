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

var CdniService = function($http, ENV, messageModel) {

	this.getCdniConfigRequests = function() {
		return $http.get(ENV.api.unstable + 'OC/CI/configuration/requests').then(
			function (result) {
				return result.data.response;
				},
			function (err) {
				messageModel.setMessages(err.data.alerts, true);
				throw err;
			}
		)
	};

	this.getCdniConfigRequestById = function(id) {
		return $http.get(ENV.api.unstable + 'OC/CI/configuration/requests?id=' + id).then(
			function (result) {
				if (result.data.response.length > 0) {
					return result.data.response[0]
				}
				return result.data.response;
			},
			function (err) {
				messageModel.setMessages(err.data.alerts, true);
				throw err;
			}
		)
	};

	this.getCurrentCdniConfigByUCDN = function(ucdn) {
		return $http.get(ENV.api.unstable + 'OC/FCI/advertisement?ucdn=' + ucdn).then(
			function (result) {
				return result.data;
			},
			function (err) {
				messageModel.setMessages(err.data.alerts, true);
				throw err;
			}
		)
	};

	this.sendResponseToCdniRequest = function(id, approve) {
		return $http.put(ENV.api.unstable + 'OC/CI/configuration/request/' + id + '/' + approve).then(
			function (result) {
				return result.data.response;
			},
			function (err) {
				messageModel.setMessages(err.data.alerts, true);
				throw err;
			}
		)
	};
};

CdniService.$inject = ['$http', 'ENV', 'messageModel'];
module.exports = CdniService;
