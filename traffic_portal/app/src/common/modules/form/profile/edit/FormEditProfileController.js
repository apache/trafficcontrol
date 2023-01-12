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
 * @param {*} profile
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/ProfileService")} profileService
 * @param {import("../../../../models/MessageModel")} messageModel
 */
var FormEditProfileController = function(profile, $scope, $controller, $uibModal, $anchorScroll, locationUtils, profileService, messageModel) {

    // extends the FormProfileController to inherit common methods
    angular.extend(this, $controller('FormProfileController', { profile: profile, $scope: $scope }));

    var deleteProfile = function(profile) {
        profileService.deleteProfile(profile.id)
            .then(function(result) {
                messageModel.setMessages(result.alerts, true);
                locationUtils.navigateToPath('/profiles');
            },
            function() {
                // do nothing
            });
    };

    $scope.profileName = angular.copy(profile.name);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(profile) {
        profileService.updateProfile(profile).
            then(function() {
                $scope.profileName = angular.copy(profile.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(profile) {
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

};

FormEditProfileController.$inject = ['profile', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'profileService', 'messageModel'];
module.exports = FormEditProfileController;
