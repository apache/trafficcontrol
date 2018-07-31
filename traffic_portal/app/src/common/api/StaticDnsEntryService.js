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

var StaticDnsEntryService = function($http, $q, Restangular, locationUtils, messageModel, ENV) {

	this.getStaticDnsEntries = function(queryParams) {
		return Restangular.all('staticdnsentries').getList(queryParams);
	};

    this.createDeliveryServiceStaticDnsEntry = function(staticDnsEntry) {
        return Restangular.service('staticdnsentries').post(staticDnsEntry)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Static DNS Entry created' } ], true);
                    locationUtils.navigateToPath('/delivery-services/' + staticDnsEntry.deliveryServiceId + '/static-dns-entries');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteDeliveryServiceStaticDnsEntry = function(queryParams) {
        return Restangular.all('staticdnsentries').remove(queryParams)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Static DNS Entry deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

    this.updateDeliveryServiceStaticDnsEntry = function(id, staticDnsEntry) {
        var request = $q.defer();

        $http.put(ENV.api['root'] + "staticdnsentries?id=" + id, staticDnsEntry)
            .then(
                function(response) {
                    messageModel.setMessages(response.data.alerts, false);
                    request.resolve();
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                    request.reject();
                }
            );
        return request.promise;
    };
};

StaticDnsEntryService.$inject = ['$http', '$q', 'Restangular', 'locationUtils', 'messageModel', 'ENV'];
module.exports = StaticDnsEntryService;
