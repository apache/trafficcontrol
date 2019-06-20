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

var PhysLocationService = function(Restangular, locationUtils, messageModel) {

    this.getPhysLocations = function(queryParams) {
        return Restangular.all('phys_locations').getList(queryParams);
    };

    this.getPhysLocation = function(id) {
        return Restangular.one("phys_locations", id).get();
    };

    this.createPhysLocation = function(physLocation) {
        return Restangular.service('phys_locations').post(physLocation)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Physical location created' } ], true);
                    locationUtils.navigateToPath('/phys-locations');

                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updatePhysLocation = function(physLocation) {
        return physLocation.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Physical location updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deletePhysLocation = function(id) {
        return Restangular.one("phys_locations", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Physical location deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

};

PhysLocationService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = PhysLocationService;
