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

var TenantService = function(Restangular, messageModel) {

    this.getTenants = function(queryParams) {
        return Restangular.all('tenants').getList(queryParams);
    };

    this.getTenant = function(id) {
        return Restangular.one("tenants", id).get();
    };

    this.createTenant = function(tenant) {
        return Restangular.service('tenants').post(tenant)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Tenant created' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
        );
    };

    this.updateTenant = function(tenant) {
        return tenant.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Tenant updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteTenant = function(id) {
        return Restangular.one("tenants", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Tenant deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

};

TenantService.$inject = ['Restangular', 'messageModel'];
module.exports = TenantService;
