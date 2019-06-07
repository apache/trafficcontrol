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

var DeliveryServiceRequestService = function(Restangular, $http, $q, locationUtils, messageModel, ENV) {

	this.getDeliveryServiceRequests = function(queryParams) {
		return Restangular.all('deliveryservice_requests').getList(queryParams);
	};

	this.createDeliveryServiceRequest = function(dsRequest) {
		var request = $q.defer();

		// strip out any falsy values or duplicates from consistentHashQueryParams
		dsRequest.deliveryService.consistentHashQueryParams = Array.from(new Set(dsRequest.deliveryService.consistentHashQueryParams)).filter(function(i){return i;});

		$http.post(ENV.api['root'] + "deliveryservice_requests", dsRequest)
			.then(
				function(result) {
					request.resolve(result.data.response);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject(fault);
				}
			);

		return request.promise;
	};

	this.updateDeliveryServiceRequest = function(id, dsRequest) {
		var request = $q.defer();

		// strip out any falsy values or duplicates from consistentHashQueryParams
		dsRequest.deliveryService.consistentHashQueryParams = Array.from(new Set(dsRequest.deliveryService.consistentHashQueryParams)).filter(function(i){return i;});

		$http.put(ENV.api['root'] + "deliveryservice_requests?id=" + id, dsRequest)
			.then(
				function() {
					request.resolve();
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject();
				}
			);

		return request.promise;
	};

	this.deleteDeliveryServiceRequest = function(id, delay) {
		var deferred = $q.defer();

		$http.delete(ENV.api['root'] + "deliveryservice_requests?id=" + id)
			.then(
				function(response) {
					deferred.resolve(response);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					deferred.reject(fault);
				}
			);

		return deferred.promise;
	};

	this.assignDeliveryServiceRequest = function(id, userId) {
		var request = $q.defer();

		$http.put(ENV.api['root'] + "deliveryservice_requests/" + id + "/assign", { id: id, assigneeId: userId })
			.then(
				function() {
					request.resolve();
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject();
				}
			);

		return request.promise;
	};

	this.updateDeliveryServiceRequestStatus = function(id, status) {
		var request = $q.defer();

		$http.put(ENV.api['root'] + "deliveryservice_requests/" + id + "/status", { id: id, status: status })
			.then(
				function() {
					request.resolve();
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject();
				}
			);

		return request.promise;
	};

	this.getDeliveryServiceRequestComments = function(queryParams) {
		return Restangular.all('deliveryservice_request_comments').getList(queryParams);
	};

	this.createDeliveryServiceRequestComment = function(comment) {
		var request = $q.defer();

		$http.post(ENV.api['root'] + "deliveryservice_request_comments", comment)
			.then(
				function(response) {
					request.resolve(response);
				},
				function(fault) {
					request.reject(fault);
				}
			);

		return request.promise;
	};

	this.updateDeliveryServiceRequestComment = function(comment) {
		var request = $q.defer();

		$http.put(ENV.api['root'] + "deliveryservice_request_comments?id=" + comment.id, comment)
			.then(
				function() {
					request.resolve();
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject();
				}
			);

		return request.promise;
	};

	this.deleteDeliveryServiceRequestComment = function(comment) {
		var deferred = $q.defer();

		$http.delete(ENV.api['root'] + "deliveryservice_request_comments?id=" + comment.id)
			.then(
				function(response) {
					deferred.resolve(response);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					deferred.reject(fault);
				}
			);

		return deferred.promise;
	};

};

DeliveryServiceRequestService.$inject = ['Restangular', '$http', '$q', 'locationUtils', 'messageModel', 'ENV'];
module.exports = DeliveryServiceRequestService;
