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

var TableProfilesController = function(profiles, $scope, $state, $location, $uibModal, locationUtils, profileService) {

    $scope.profiles = profiles;

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
            title: 'Compare Profile Parameters',
            message: "Please select 2 profiles to compare parameters"
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
            $location.path($location.path() + '/compare/' + profiles[0].id + '/' + profiles[1].id);
        }, function () {
            // do nothing
        });
    };


    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    angular.element(document).ready(function () {
        $('#profilesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "aaSorting": []
        });
    });

};

TableProfilesController.$inject = ['profiles', '$scope', '$state', '$location', '$uibModal', 'locationUtils', 'profileService'];
module.exports = TableProfilesController;
