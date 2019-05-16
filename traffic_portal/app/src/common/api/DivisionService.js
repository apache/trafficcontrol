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

var DivisionService = function(Restangular, locationUtils, messageModel) {

    this.getDivisions = function(queryParams) {
        return Restangular.all('divisions').getList(queryParams);
    };

    this.getDivision = function(id) {
        return Restangular.one("divisions", id).get();
    };

    this.createDivision = function(division) {
        return Restangular.service('divisions').post(division)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Division created' } ], true);
                    locationUtils.navigateToPath('/divisions');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateDivision = function(division) {
        return division.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Division updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteDivision = function(id) {
        return Restangular.one("divisions", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Division deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

};

DivisionService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = DivisionService;
