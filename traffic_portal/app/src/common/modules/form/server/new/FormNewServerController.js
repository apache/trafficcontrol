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

var FormNewServerController = function(server, $scope, $controller, serverService, statusService) {

    // extends the FormServerController to inherit common methods
    angular.extend(this, $controller('FormServerController', { server: server, $scope: $scope }));

    var getStatuses = function() {
        statusService.getStatuses()
            .then(function(result) {
                $scope.statuses = result;
                // Issue #2651 - Enabling server status for New Server but still defaulting enabled dropdown to OFFLINE
                var offlineStatus = _.find(result, function(status){ return status.name == 'OFFLINE' });
                $scope.server.statusId = offlineStatus.id;
            });
    };

    $scope.serverName = 'New';

    $scope.settings = {
        isNew: true,
        saveLabel: 'Create'
    };

    $scope.save = function(server) {
        serverService.createServer(server);
    };

    var init = function () {
        getStatuses();
    };
    init();

};

FormNewServerController.$inject = ['server', '$scope', '$controller', 'serverService', 'statusService'];
module.exports = FormNewServerController;
