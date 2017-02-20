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

var DeliveryServiceRegexService = function(Restangular, locationUtils, messageModel) {

	this.getDeliveryServiceRegexes = function(dsId) {
		return Restangular.one('deliveryservices', dsId).getList('regexes');
	};

	this.getDeliveryServiceRegex = function(dsId, regexId) {
		return Restangular.one('deliveryservices', dsId).one('regexes', regexId).get();
	};

	this.createDeliveryServiceRegex = function(dsId, regex) {
		return Restangular.one('deliveryservices', dsId).all('regexes').post(regex)
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Regex created' } ], true);
					locationUtils.navigateToPath('/configure/delivery-services/' + dsId + '/regexes');
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

	this.updateDeliveryServiceRegex = function(regex) {
		return regex.put()
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Regex updated' } ], false);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

	this.deleteDeliveryServiceRegex = function(dsId, regexId) {
		return Restangular.one('deliveryservices', dsId).one('regexes', regexId).remove()
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Regex deleted' } ], true);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, true);
				}
			);
	};

};

DeliveryServiceRegexService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = DeliveryServiceRegexService;
