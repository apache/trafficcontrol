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

var FormEditParameterController = function(parameter, $scope, $controller, $uibModal, $anchorScroll, locationUtils, parameterService) {

    // extends the FormParameterController to inherit common methods
    angular.extend(this, $controller('FormParameterController', { parameter: parameter, $scope: $scope }));

    var deleteParameter = function(parameter) {
        parameterService.deleteParameter(parameter.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/parameters');
            });
    };

    $scope.parameterName = angular.copy(parameter.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(parameter) {
        parameterService.updateParameter(parameter).
            then(function() {
                $scope.parameterName = angular.copy(parameter.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(parameter) {
        var params = {
            title: 'Delete Parameter: ' + parameter.name,
            key: parameter.name
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
            deleteParameter(parameter);
        }, function () {
            // do nothing
        });
    };

};

FormEditParameterController.$inject = ['parameter', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'parameterService'];
module.exports = FormEditParameterController;
