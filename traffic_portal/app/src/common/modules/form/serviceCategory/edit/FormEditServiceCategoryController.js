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
 * @param {*} serviceCategory
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/ServiceCategoryService")} serviceCategoryService
 */
var FormEditServiceCategoryController = function(serviceCategory, $scope, $controller, $uibModal, $anchorScroll, locationUtils, serviceCategoryService) {

    // extends the FormServiceCategoryController to inherit common methods
    angular.extend(this, $controller('FormServiceCategoryController', { serviceCategory: serviceCategory, $scope: $scope }));

    var deleteServiceCategory = function(serviceCategory) {
        serviceCategoryService.deleteServiceCategory(serviceCategory.name)
            .then(function() {
                locationUtils.navigateToPath('/service-categories');
            });
    };

    $scope.serviceCategoryName = angular.copy(serviceCategory.name);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(serviceCategory) {
        serviceCategoryService.updateServiceCategory(serviceCategory, $scope.serviceCategoryName).
            then(function() {
                $scope.serviceCategoryName = angular.copy(serviceCategory.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(serviceCategory) {
        var params = {
            title: 'Delete Service Category: ' + serviceCategory.name,
            key: serviceCategory.name
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
            deleteServiceCategory(serviceCategory);
        }, function () {
            // do nothing
        });
    };

};

FormEditServiceCategoryController.$inject = ['serviceCategory', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'serviceCategoryService'];
module.exports = FormEditServiceCategoryController;
