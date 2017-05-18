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

var FormServerController = function(server, $scope, $location, $state, $uibModal, formUtils, locationUtils, serverService, cacheGroupService, cdnService, physLocationService, profileService, statusService, typeService) {

    var getPhysLocations = function() {
        physLocationService.getPhysLocations()
            .then(function(result) {
                $scope.physLocations = result;
            });
    };

    var getCacheGroups = function() {
        cacheGroupService.getCacheGroups()
            .then(function(result) {
                $scope.cacheGroups = result;
            });
    };

    var getTypes = function() {
        typeService.getTypes({ useInTable: 'server' })
            .then(function(result) {
                $scope.types = result;
            });
    };

    var getCDNs = function() {
        cdnService.getCDNs()
            .then(function(result) {
                $scope.cdns = result;
            });
    };

    var getStatuses = function() {
        statusService.getStatuses()
            .then(function(result) {
                $scope.statuses = result;
            });
    };

    var getProfiles = function() {
        profileService.getProfiles()
            .then(function(result) {
                $scope.profiles = result;
            });
    };

    var updateStatus = function(status) {
        // todo: hook this into PUT /api/version/server/:id/status
        console.log(status.name);
        console.log(status.offlineReason);
        alert('this still needs to be hooked into the api');
    };
    
    var refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };


    // supposedly matches IPv4 and IPv6 formats. but actually need one that matches each. todo.
    $scope.validations = {
        ipRegex: new RegExp(/^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z]|[A-Za-z][A-Za-z0-9\-]*[A-Za-z0-9])$|^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$/)
    };


    $scope.server = server;

    $scope.falseTrue = [
        { value: false, label: 'false' },
        { value: true, label: 'true' }
    ];

    $scope.queueServerUpdates = function(server) {
        serverService.queueServerUpdates(server.id)
            .then(
                function() {
                    refresh();
                }
            );
    };

    $scope.clearServerUpdates = function(server) {
        serverService.clearServerUpdates(server.id)
            .then(
                function() {
                    refresh();
                }
            );
    };

    $scope.confirmStatusUpdate = function() {
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/status/dialog.select.status.tpl.html',
            controller: 'DialogSelectStatusController',
            size: 'md',
            resolve: {
                statuses: function() {
                    return $scope.statuses;
                }
            }
        });
        modalInstance.result.then(function(status) {
            updateStatus(status);
        }, function () {
            // do nothing
        });
    };


    $scope.viewConfig = function() {
        alert('not hooked up yet: view config files for server');
    };

    $scope.viewDeliveryServices = function() {
        $location.path($location.path() + '/delivery-services');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getPhysLocations();
        getCacheGroups();
        getTypes();
        getCDNs();
        getStatuses();
        getProfiles();
    };
    init();

};

FormServerController.$inject = ['server', '$scope', '$location', '$state', '$uibModal', 'formUtils', 'locationUtils', 'serverService', 'cacheGroupService', 'cdnService', 'physLocationService', 'profileService', 'statusService', 'typeService'];
module.exports = FormServerController;
