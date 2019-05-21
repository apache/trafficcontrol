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

var TypeService = function(Restangular, locationUtils, messageModel) {

    this.getTypes = function(queryParams) {
        return Restangular.all('types').getList(queryParams);
    };

    this.getType = function(id) {
        return Restangular.one("types", id).get();
    };

    this.createType = function(type) {
        return Restangular.service('types').post(type)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Type created' } ], true);
                    locationUtils.navigateToPath('/types');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateType = function(type) {
        return type.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Type updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
        );
    };

    this.deleteType = function(id) {
        return Restangular.one("types", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Type deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
        );
    };

};

TypeService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = TypeService;
