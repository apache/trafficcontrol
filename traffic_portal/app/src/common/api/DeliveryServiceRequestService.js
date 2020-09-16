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

var DeliveryServiceRequestService = function($http, locationUtils, messageModel, ENV) {

	this.getDeliveryServiceRequests = function(queryParams) {
		return $http.get(ENV.api['root'] + 'deliveryservice_requests', {params: queryParams}).then(
			function(result) {
				return result.data.response;
			},
			function(err) {
				throw err;
			}
		);
	};

	this.createDeliveryServiceRequest = function(dsRequest) {

		// strip out any falsy values or duplicates from consistentHashQueryParams
		dsRequest.deliveryService.consistentHashQueryParams = Array.from(new Set(dsRequest.deliveryService.consistentHashQueryParams)).filter(function(i){return i;});

		return $http.post(ENV.api['root'] + "deliveryservice_requests", dsRequest).then(
			function(result) {
				return result.data.response;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.updateDeliveryServiceRequest = function(id, dsRequest) {

		// strip out any falsy values or duplicates from consistentHashQueryParams
		dsRequest.deliveryService.consistentHashQueryParams = Array.from(new Set(dsRequest.deliveryService.consistentHashQueryParams)).filter(function(i){return i;});

		return $http.put(ENV.api['root'] + "deliveryservice_requests", dsRequest, {params: {id: id}}).then(
			function(result) {
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.deleteDeliveryServiceRequest = function(id, delay) {
		return $http.delete(ENV.api['root'] + "deliveryservice_requests", {params: {id: id}}).then(
			function(response) {
				return response;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.assignDeliveryServiceRequest = function(id, userId) {
		return $http.put(ENV.api['root'] + "deliveryservice_requests/" + id + "/assign", { id: id, assigneeId: userId }).then(
			function(result) {
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.updateDeliveryServiceRequestStatus = function(id, status) {
		return $http.put(ENV.api['root'] + "deliveryservice_requests/" + id + "/status", { id: id, status: status }).then(
			function(result) {
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.getDeliveryServiceRequestComments = function(queryParams) {
		return $http.get(ENV.api['root'] + 'deliveryservice_request_comments', {params: queryParams}).then(
			function(result) {
				return result.data.response;
			},
			function(err) {
				throw err;
			}
		);
	};

	this.createDeliveryServiceRequestComment = function(comment) {
		return $http.post(ENV.api['root'] + "deliveryservice_request_comments", comment).then(
			function(response) {
				return response;
			},
			function(err) {
				throw err;
			}
		);
	};

	this.updateDeliveryServiceRequestComment = function(comment) {
		return $http.put(ENV.api['root'] + "deliveryservice_request_comments", comment, {params: {id: comment.id}}).then(
				function(result) {
					return result;
				},
				function(err) {
					messageModel.setMessages(err.data.alerts, false);
					throw err;
				}
			);
	};

	this.deleteDeliveryServiceRequestComment = function(comment) {
		return $http.delete(ENV.api['root'] + "deliveryservice_request_comments", {params: {id: comment.id}}).then(
			function(response) {
				return response;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};
};

DeliveryServiceRequestService.$inject = ['$http', 'locationUtils', 'messageModel', 'ENV'];
module.exports = DeliveryServiceRequestService;
