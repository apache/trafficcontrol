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
 * @param {*} origin
 * @param {*} $scope
 * @param {*} $state
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/OriginService")} originService
 */
var FormEditOriginController = function(origin, $scope, $state, $controller, $uibModal, $anchorScroll, locationUtils, originService) {

    $scope.origin = origin[0]

    // extends the FormOriginController to inherit common methods
    angular.extend(this, $controller('FormOriginController', { origin: $scope.origin, $scope: $scope }));

    $scope.originName = angular.copy($scope.origin.name);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update',
        deleteLabel: 'Delete'
    };

    $scope.save = function(origin) {
        originService.updateOrigin(origin.id, origin).
            then(
                function(result) {
                    $state.reload(); // reloads all the resolves for the view
                },
                function(fault) {
                    $anchorScroll(); // scrolls window to top
                }
            );
    };

    $scope.confirmDelete = function(origin) {
        var params = {
            title: 'Delete Origin: ' + origin.name,
            key: origin.name
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
            originService.deleteOrigin(origin.id)
                .then(
                    function(result) {
                        locationUtils.navigateToPath('/origins');
                    },
                    function(fault) {
                        $anchorScroll(); // scrolls window to top
                    }
                );
        }, function () {
            // do nothing
        });
    };

};

FormEditOriginController.$inject = ['origin', '$scope', '$state', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'originService'];
module.exports = FormEditOriginController;
