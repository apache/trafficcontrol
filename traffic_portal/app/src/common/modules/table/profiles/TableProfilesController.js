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

var TableProfilesController = function(profiles, $scope, $state, $location, $uibModal, $window, locationUtils, profileService, messageModel, fileUtils) {

    var confirmDelete = function(profile) {
        var params = {
            title: 'Delete Profile: ' + profile.name,
            key: profile.name
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
            deleteProfile(profile);
        }, function () {
            // do nothing
        });
    };

    var deleteProfile = function(profile) {
        profileService.deleteProfile(profile.id)
            .then(function(result) {
                messageModel.setMessages(result.alerts, false);
                $scope.refresh();
            });
    };

    var cloneProfile = function(profile) {
        var params = {
            title: 'Clone Profile',
            message: "Your are about to clone the " + profile.name + " profile. Your clone will have the same attributes and parameter assignments as the " + profile.name + " profile.<br><br>Please enter a name for your cloned profile."
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/input/dialog.input.tpl.html',
            controller: 'DialogInputController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function(clonedProfileName) {
            profileService.cloneProfile(profile.name, clonedProfileName);
        }, function () {
            // do nothing
        });
    };

    var exportProfile = function(profile) {
        profileService.exportProfile(profile.id).
        then(
            function(result) {
                fileUtils.exportJSON(result, profile.name, 'traffic_ops');
            }
        );

    };

    $scope.profiles = profiles;

    $scope.contextMenuItems = [
        {
            text: 'Open in New Tab',
            click: function ($itemScope) {
                $window.open('/#!/profiles/' + $itemScope.p.id, '_blank');
            }
        },
        null, // Divider
        {
            text: 'Edit',
            click: function ($itemScope) {
                $scope.editProfile($itemScope.p.id);
            }
        },
        {
            text: 'Delete',
            click: function ($itemScope) {
                confirmDelete($itemScope.p);
            }
        },
        null, // Divider
        {
            text: 'Clone Profile',
            click: function ($itemScope) {
                cloneProfile($itemScope.p);
            }
        },
        {
            text: 'Export Profile',
            click: function ($itemScope) {
                exportProfile($itemScope.p);
            }
        },
        null, // Divider
        {
            text: 'Manage Parameters',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/profiles/' + $itemScope.p.id + '/parameters');
            }
        },
        {
            text: 'Manage Servers',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/profiles/' + $itemScope.p.id + '/servers');
            }
        }
    ];

    $scope.editProfile = function(id) {
        locationUtils.navigateToPath('/profiles/' + id);
    };

    $scope.createProfile = function() {
        locationUtils.navigateToPath('/profiles/new');
    };

    $scope.importProfile = function() {
        var params = {
            title: 'Import Profile',
            message: "Drop Profile Here"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/import/dialog.import.tpl.html',
            controller: 'DialogImportController',
            size: 'lg',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function(importJSON) {
            profileService.importProfile(importJSON);
        }, function () {
            // do nothing
        });
    };

    $scope.compareProfiles = function() {
        var params = {
            title: 'Compare Profiles',
            message: 'Please select 2 profiles to compare',
            labelFunction: function(item) { return item['name'] + ' (' + item['type'] + ')' }
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/compare/dialog.compare.tpl.html',
            controller: 'DialogCompareController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function(profileService) {
                    return profileService.getProfiles({ orderby: 'name' });
                }
            }
        });
        modalInstance.result.then(function(profiles) {
            $location.path($location.path() + '/' + profiles[0].id + '/' + profiles[1].id + '/compare/diff');
        }, function () {
            // do nothing
        });
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    angular.element(document).ready(function () {
        $('#profilesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "aaSorting": []
        });
    });

};

TableProfilesController.$inject = ['profiles', '$scope', '$state', '$location', '$uibModal', '$window', 'locationUtils', 'profileService', 'messageModel', 'fileUtils'];
module.exports = TableProfilesController;
