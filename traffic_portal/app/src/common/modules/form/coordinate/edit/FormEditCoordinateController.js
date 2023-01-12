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
 *
 * @param {*} coordinate
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/CoordinateService")} coordinateService
 */
var FormEditCoordinateController = function(coordinate, $scope, $controller, $uibModal, $anchorScroll, locationUtils, coordinateService) {

    $scope.coordinate = coordinate[0]

    // extends the FormCoordinateController to inherit common methods
    angular.extend(this, $controller('FormCoordinateController', { coordinate: $scope.coordinate, $scope: $scope }));

    $scope.coordinateName = angular.copy($scope.coordinate.name);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    var deleteCoordinate = function(coordinate) {
        coordinateService.deleteCoordinate(coordinate.id)
            .then(function() {
                locationUtils.navigateToPath('/coordinates');
            });
    };

    $scope.save = function(coordinate) {
        coordinateService.updateCoordinate(coordinate.id, coordinate).
            then(function() {
                $scope.coordinateName = angular.copy(coordinate.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(coordinate) {
        var params = {
            title: 'Delete Coordinate: ' + coordinate.name,
            key: coordinate.name
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
            deleteCoordinate(coordinate);
        }, function () {
            // do nothing
        });
    };
};

FormEditCoordinateController.$inject = ['coordinate', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'coordinateService'];
module.exports = FormEditCoordinateController;
