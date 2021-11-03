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

var TenantService = function($http, ENV, messageModel) {

    this.getTenants = function(queryParams) {
        return $http.get(ENV.api.unstable + 'tenants', {params: queryParams}).then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        )
    };

    this.getTenant = function(id) {
        return $http.get(ENV.api.unstable + 'tenants', {params: {id: id}}).then(
            function (result) {
                return result.data.response[0];
            },
            function (err) {
                throw err;
            }
        )
    };

    this.createTenant = function(tenant) {
        return $http.post(ENV.api.unstable + 'tenants', tenant).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Tenant created' } ], true);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.updateTenant = function(tenant) {
        return $http.put(ENV.api.unstable + 'tenants/' + tenant.id, tenant).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Tenant updated' } ], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.deleteTenant = function(id) {
        return $http.delete(ENV.api.unstable + "tenants/" + id).then(
            function(result) {
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

};

TenantService.$inject = ['$http', 'ENV', 'messageModel'];
module.exports = TenantService;
