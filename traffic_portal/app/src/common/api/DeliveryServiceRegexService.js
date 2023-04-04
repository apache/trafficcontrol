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

/**
 * @param {import("angular").IHttpService} $http
 * @param {import("../service/utils/LocationUtils")} locationUtils
 * @param {import("../models/MessageModel")} messageModel
 * @param {{api: Record<PropertyKey, string>}} ENV
 */
var DeliveryServiceRegexService = function($http, locationUtils, messageModel, ENV) {

	this.getDeliveryServiceRegexes = function(dsId) {
		return $http.get(ENV.api.unstable + 'deliveryservices/' + dsId + '/regexes').then(
			function(result) {
				return result.data.response;
			},
			function(err) {
				throw err;
			}
		);
	};

	this.getDeliveryServiceRegex = function(dsId, regexId) {
		return $http.get(ENV.api.unstable + 'deliveryservices/' + dsId + '/regexes', {params: {id: regexId}}).then(
			function(result) {
				return result.data.response[0];
			},
			function(err) {
				throw err;
			}
		);
	};

	this.createDeliveryServiceRegex = function(dsId, regex) {
		return $http.post(ENV.api.unstable + 'deliveryservices/' + dsId + '/regexes', regex).then(
			function(result) {
				messageModel.setMessages(result.data.alerts, true);
				locationUtils.navigateToPath('/delivery-services/' + dsId + '/regexes');
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.updateDeliveryServiceRegex = function(dsId, regex) {
		return $http.put(ENV.api.unstable + 'deliveryservices/' + dsId + '/regexes/' + regex.id, regex).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text:'Regex updated'}], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.deleteDeliveryServiceRegex = function(dsId, regexId) {
		return $http.delete(ENV.api.unstable + 'deliveryservices/' + dsId + '/regexes/' + regexId).then(
			function(result) {
				messageModel.setMessages(result.data.alerts, true);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

};

DeliveryServiceRegexService.$inject = ['$http', 'locationUtils', 'messageModel', 'ENV'];
module.exports = DeliveryServiceRegexService;
