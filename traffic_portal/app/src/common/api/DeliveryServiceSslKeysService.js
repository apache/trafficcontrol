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

var DeliveryServiceSslKeysService = function($http, $q, locationUtils, messageModel, ENV) {
	this.generateSslKeys = function(deliveryService, sslKeys, generateSslKeyForm) {
		if (sslKeys.hasOwnProperty('version')){
			generateSslKeyForm.version = parseInt(sslKeys.version) + 1;
		} else {
			generateSslKeyForm.version = 1;
		}

		generateSslKeyForm.cdn = deliveryService.cdnName;
		generateSslKeyForm.deliveryservice = deliveryService.xmlId;
		generateSslKeyForm.key = deliveryService.xmlId;

		var request = $q.defer();
        $http.post(ENV.api['root'] + "deliveryservices/sslkeys/generate", generateSslKeyForm)
        .then(
            function(result) {
            	messageModel.setMessages([ { level: 'success', text: 'SSL Keys generated and updated for ' + deliveryService.xmlId } ], true);
                request.resolve(result.data.response);
            },
            function(fault) {
            	messageModel.setMessages(fault.data.alerts, false);
                request.reject(fault);
            }
        );
        return request.promise;
	};

	this.addSslKeys = function(sslKeys, deliveryService) {
		var request = $q.defer();

        sslKeys.key = deliveryService.xmlId;
        if (sslKeys.hasOwnProperty('version')){
            sslKeys.version = parseInt(sslKeys.version) + 1;
        } else {
            sslKeys.version = 1;
        }

        sslKeys.cdn = deliveryService.cdnName;
        sslKeys.deliveryservice = deliveryService.xmlId;

        $http.post(ENV.api['root'] + "deliveryservices/sslkeys/add", sslKeys)
        .then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
                request.resolve(result.data.response);
            },
            function(fault) {
            	messageModel.setMessages(fault.data.alerts, false);
                request.reject(fault);
            }
        );
        return request.promise;
	};

	this.getSslKeys = function(deliveryService) {
		var request = $q.defer();
        $http.get(ENV.api['root'] + "deliveryservices/xmlId/" + deliveryService.xmlId + "/sslkeys?decode=true")
        .then(
            function(result) {
                request.resolve(result.data.response);
            },
            function(fault) {
            	messageModel.setMessages(fault.data.alerts, true);
                request.reject(fault);
            }
        );
        return request.promise;
	};
};

DeliveryServiceSslKeysService.$inject = ['$http', '$q', 'locationUtils', 'messageModel', 'ENV'];
module.exports = DeliveryServiceSslKeysService;
