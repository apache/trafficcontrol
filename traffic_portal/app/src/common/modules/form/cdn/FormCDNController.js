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

function FormCDNController(cdn, $scope, $uibModal, formUtils, cdnService) {

    function queueServerUpdates(cdn) {
        cdnService.queueServerUpdates(cdn.id);
    }

    function clearServerUpdates(cdn) {
        cdnService.clearServerUpdates(cdn.id);
    }

    $scope.cdn = cdn;

    $scope.queueServerUpdates = function(cdn) {
        var params = {
            title: 'Queue Server Updates: ' + cdn.name,
            message: 'Are you sure you want to queue server updates for all ' + cdn.name + ' servers?'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            queueServerUpdates(cdn);
        }, function () {
            // do nothing
        });
    };

    $scope.clearServerUpdates = function(cdn) {
        var params = {
            title: 'Clear Server Updates: ' + cdn.name,
            message: 'Are you sure you want to clear server updates for all ' + cdn.name + ' servers?'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            clearServerUpdates(cdn);
        }, function () {
            // do nothing
        });
    };

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

}

FormCDNController.$inject = ['cdn', '$scope', '$uibModal', 'formUtils', 'cdnService'];
module.exports = FormCDNController;
