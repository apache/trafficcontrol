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

var CapabilityService = function(Restangular, $q, $http, messageModel, ENV) {

	this.getCapabilities = function(queryParams) {
		return Restangular.all('capabilities').getList(queryParams);
	};

	this.getCapability = function(name) {
		return Restangular.one("capabilities", name).get();
	};

	this.createCapability = function(cap) {
		var request = $q.defer();

		$http.post(ENV.api['root'] + "capabilities", cap)
			.then(
				function(result) {
					request.resolve(result.data);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject(fault);
				}
			);

		return request.promise;
	};

	this.updateCapability = function(cap) {
		var request = $q.defer();

		$http.put(ENV.api['root'] + "capabilities/" + cap.name, cap)
			.then(
				function(result) {
					request.resolve(result.data);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject();
				}
			);

		return request.promise;
	};

	this.deleteCapability = function(cap) {
		var request = $q.defer();

		$http.delete(ENV.api['root'] + "capabilities/" + cap.name)
			.then(
				function(result) {
					request.resolve(result.data);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject(fault);
				}
			);

		return request.promise;
	};

};

CapabilityService.$inject = ['Restangular', '$q', '$http', 'messageModel', 'ENV'];
module.exports = CapabilityService;
