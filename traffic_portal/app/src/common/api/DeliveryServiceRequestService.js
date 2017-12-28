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

var DeliveryServiceRequestService = function(Restangular, locationUtils, messageModel, dsModel) {

	this.getDeliveryServiceRequests = function(queryParams) {
		// return Restangular.all('deliveryservice_requests').getList(queryParams);
		return dsModel.requests[0];
	};

	this.createDeliveryServiceRequest = function(dsRequest) {
		dsModel.requests = [ dsRequest ];
		locationUtils.navigateToPath('/delivery-service-requests');

		// return Restangular.service('deliveryservice_requests').post(dsRequest)
		// 	.then(
		// 		function() {
		// 			messageModel.setMessages([ { level: 'success', text: 'DS request created' } ], true);
		// 			locationUtils.navigateToPath('/delivery-services');
		// 		},
		// 		function(fault) {
		// 			messageModel.setMessages(fault.data.alerts, false);
		// 		}
		// 	);
	};

	this.updateDeliveryServiceRequest = function(dsRequest) {
		dsModel.requests = [ dsRequest ];
		locationUtils.navigateToPath('/delivery-service-requests');
	};


	this.deleteDeliveryServiceRequest = function(dsRequestId) {
		dsModel.requests = [];
		locationUtils.navigateToPath('/delivery-service-requests');
	};

};

DeliveryServiceRequestService.$inject = ['Restangular', 'locationUtils', 'messageModel', 'dsModel'];
module.exports = DeliveryServiceRequestService;
