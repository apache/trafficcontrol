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

	this.createDeliveryServiceRequest = function(dsRequest, delay) {
		return Restangular.service('deliveryservice_requests').post(dsRequest)
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Created request to ' + dsRequest.changeType + ' the ' + dsRequest.request.xmlId + ' delivery service' } ], delay);
					locationUtils.navigateToPath('/delivery-service-requests');
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

	this.updateDeliveryServiceRequest = function(id, dsRequest) {
		var request = $q.defer();

		$http.put(ENV.api['root'] + "deliveryservice_requests/" + id, dsRequest)
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
		return Restangular.one("deliveryservice_requests", id).remove()
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Delivery service request deleted' } ], delay);
					locationUtils.navigateToPath('/delivery-service-requests');
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

};

DeliveryServiceRequestService.$inject = ['Restangular', '$http', '$q', 'locationUtils', 'messageModel', 'ENV'];
module.exports = DeliveryServiceRequestService;
