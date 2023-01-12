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

/**
 * @param {*} server
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/ServerService")} serverService
 * @param {import("../../../../api/StatusService")} statusService
 * @param {import("../../../../models/MessageModel")} messageModel
 */
var FormEditServerController = function(server, $anchorScroll, $scope, $controller, $uibModal, locationUtils, serverService, statusService, messageModel) {

    // extends the FormServerController to inherit common methods
    angular.extend(this, $controller('FormServerController', { server: server[0], $scope: $scope }));

    var getStatuses = function() {
        statusService.getStatuses()
            .then(function(result) {
                $scope.statuses = result;
            });
    };

    var deleteServer = function(server) {
        serverService.deleteServer(server.id)
            .then(function(result) {
                messageModel.setMessages(result.alerts, true);
                locationUtils.navigateToPath('/servers');
            });
    };

    $scope.serverName = server[0].hostName;

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(server) {
        serverService.updateServer(server).
            then(
                function(result) {
                    $scope.refresh();
                    messageModel.setMessages(result.data.alerts, false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            )
            .finally(
                function() {
                    $anchorScroll(); // scrolls window to top for message
                }
            );
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

    var init = function () {
        getStatuses();
    };
    init();

};

FormEditServerController.$inject = ['server', '$anchorScroll', '$scope', '$controller', '$uibModal', 'locationUtils', 'serverService', 'statusService', 'messageModel'];
module.exports = FormEditServerController;
