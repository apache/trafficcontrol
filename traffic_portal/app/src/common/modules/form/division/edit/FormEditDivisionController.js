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

var FormEditDivisionController = function(division, $scope, $controller, $uibModal, $anchorScroll, locationUtils, divisionService) {

    // extends the FormDivisionController to inherit common methods
    angular.extend(this, $controller('FormDivisionController', { division: division, $scope: $scope }));

    var deleteDivision = function(division) {
        divisionService.deleteDivision(division.id)
            .then(function() {
                locationUtils.navigateToPath('/divisions');
            },
            function() {
                // Do nothing
            });
    };

    $scope.divisionName = angular.copy(division.name);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(division) {
        divisionService.updateDivision(division).
            then(function() {
                $scope.divisionName = angular.copy(division.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(division) {
        var params = {
            title: 'Delete Division: ' + division.name,
            key: division.name
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
            deleteDivision(division);
        }, function () {
            // do nothing
        });
    };

};

FormEditDivisionController.$inject = ['division', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'divisionService'];
module.exports = FormEditDivisionController;
