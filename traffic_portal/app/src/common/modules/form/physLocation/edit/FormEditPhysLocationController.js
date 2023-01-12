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
 * @param {*} physLocation
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/PhysLocationService")} physLocationService
 */
var FormEditPhysLocationController = function(physLocation, $scope, $controller, $uibModal, $anchorScroll, locationUtils, physLocationService) {

    // extends the FormPhysLocationController to inherit common methods
    angular.extend(this, $controller('FormPhysLocationController', { physLocation: physLocation, $scope: $scope }));

    var deletePhysLocation = function(physLocation) {
        physLocationService.deletePhysLocation(physLocation.id)
            .then(function() {
                locationUtils.navigateToPath('/phys-locations');
            });
    };

    $scope.physLocationName = angular.copy(physLocation.name);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(physLocation) {
        physLocationService.updatePhysLocation(physLocation).
            then(function() {
                $scope.physLocationName = angular.copy(physLocation.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(physLocation) {
        var params = {
            title: 'Delete Physical Location: ' + physLocation.name,
            key: physLocation.name
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
            deletePhysLocation(physLocation);
        }, function () {
            // do nothing
        });
    };

};

FormEditPhysLocationController.$inject = ['physLocation', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'physLocationService'];
module.exports = FormEditPhysLocationController;
