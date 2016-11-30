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

var FormEditServerController = function(server, $scope, $controller, $uibModal, $anchorScroll, locationUtils, serverService) {

    // extends the FormServerController to inherit common methods
    angular.extend(this, $controller('FormServerController', { server: server, $scope: $scope }));

    var deleteServer = function(server) {
        serverService.deleteServer(server.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/servers');
            });
    };

    $scope.serverName = angular.copy(server.hostName);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(server) {
        serverService.updateServer(server).
            then(function() {
                $scope.serverName = angular.copy(server.hostName);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(server) {
        var params = {
            title: 'Delete Server: ' + server.hostName,
            key: server.hostName
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
            controller: 'DialogDeleteController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deleteServer(server);
        }, function () {
            // do nothing
        });
    };

};

FormEditServerController.$inject = ['server', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'serverService'];
module.exports = FormEditServerController;
