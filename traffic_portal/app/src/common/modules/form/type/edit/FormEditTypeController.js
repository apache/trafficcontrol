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
 * @param {*} type
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/TypeService")} typeService
 */
var FormEditTypeController = function(type, $scope, $controller, $uibModal, $anchorScroll, locationUtils, typeService) {

    // extends the FormTypeController to inherit common methods
    angular.extend(this, $controller('FormTypeController', { type: type, $scope: $scope }));

    var deleteType = function(type) {
        typeService.deleteType(type.id)
            .then(function() {
                locationUtils.navigateToPath('/types');
            });
    };

    $scope.typeName = angular.copy(type.name);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(type) {
        typeService.updateType(type).
            then(function() {
                $scope.typeName = angular.copy(type.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(type) {
        var params = {
            title: 'Delete Type: ' + type.name,
            key: type.name
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
            deleteType(type);
        }, function () {
            // do nothing
        });
    };

};

FormEditTypeController.$inject = ['type', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'typeService'];
module.exports = FormEditTypeController;
