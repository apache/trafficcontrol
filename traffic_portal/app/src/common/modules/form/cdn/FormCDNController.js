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
 * @param {*} cdn
 * @param {*} $scope
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/CDNService")} cdnService
 */
var FormCDNController = function(cdn, $scope, $location, $uibModal, formUtils, locationUtils, cdnService) {

    var queueServerUpdates = function(cdn) {
        cdnService.queueServerUpdates(cdn.id);
    };

    var clearServerUpdates = function(cdn) {
        cdnService.clearServerUpdates(cdn.id);
    };

    $scope.cdn = cdn;

    $scope.falseTrue = [
        { value: true, label: 'true' },
        { value: false, label: 'false' }
    ];

    $scope.manageDNSSEC = function() {
        $location.path($location.path() + '/dnssec-keys');
    };

    $scope.manageFederations = function() {
        $location.path($location.path() + '/federations');
    };

    $scope.viewConfig = function() {
        $location.path($location.path() + '/config/changes');
    };

    $scope.viewProfiles = function() {
        $location.url(`/profiles?cdnName=${encodeURIComponent(cdn.name)}`);
    };

    $scope.viewServers = function() {
        $location.path($location.path() + '/servers');
    };

    $scope.viewDeliveryServices = function() {
        $location.path($location.path() + '/delivery-services');
    };

    $scope.viewNotifications = function() {
        $location.path($location.path() + '/notifications');
    };

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

    $scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormCDNController.$inject = ['cdn', '$scope', '$location', '$uibModal', 'formUtils', 'locationUtils', 'cdnService'];
module.exports = FormCDNController;
