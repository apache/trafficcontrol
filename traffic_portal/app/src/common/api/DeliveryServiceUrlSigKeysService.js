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

var DeliveryServiceUrlSigKeysService = function(locationUtils, messageModel, $http, $q, ENV) {

	this.generateUrlSigKeys = function(dsXmlId) {
		var request = $q.defer();
		$http.post(ENV.api['root'] + 'deliveryservices/xmlId/' + dsXmlId + '/urlkeys/generate')
		.then(
			function(result) {
				messageModel.setMessages(result.data.alerts, false);
				request.resolve();
			},
			function(fault) {
				messageModel.setMessages(fault.data.alerts, false);
				request.reject();
			}
		);
		return request.promise;
	};

	this.copyUrlSigKeys = function(dsXmlId, copyFromXmlId) {
		var request = $q.defer();
		 $http.post(ENV.api['root'] + 'deliveryservices/xmlId/' + dsXmlId + '/urlkeys/copyFromXmlId/' + copyFromXmlId)
		.then(
			function(result) {
				messageModel.setMessages(result.data.alerts, false);
				request.resolve();
			},
			function(fault) {
				messageModel.setMessages(fault.data.alerts, false);
				request.reject();
			}
		);
		return request.promise;
	};

	this.getDeliveryServiceUrlSigKeys = function(dsId) {
		var request = $q.defer();
        $http.get(ENV.api['root'] + "deliveryservices/" + dsId + "/urlkeys")
        .then(
            function(result) {
                request.resolve(result.data.response);
            },
            function(fault) {
            	messageModel.setMessages(fault.data.alerts, false);
                request.reject();
            }
        );
        return request.promise;
	};
};

DeliveryServiceUrlSigKeysService.$inject = ['locationUtils', 'messageModel', '$http', '$q', 'ENV'];
module.exports = DeliveryServiceUrlSigKeysService;
