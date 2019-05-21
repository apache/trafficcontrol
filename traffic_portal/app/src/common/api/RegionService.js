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

var RegionService = function(Restangular, messageModel) {

    this.getRegions = function(queryParams) {
        return Restangular.all('regions').getList(queryParams);
    };

    this.getRegion = function(id) {
        return Restangular.one("regions", id).get();
    };

    this.createRegion = function(region) {
        return Restangular.service('regions').post(region)
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Region created' } ], true);
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.updateRegion = function(region) {
        return region.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Region updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteRegion = function(id) {
        return Restangular.one("regions", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Region deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

};

RegionService.$inject = ['Restangular', 'messageModel'];
module.exports = RegionService;
