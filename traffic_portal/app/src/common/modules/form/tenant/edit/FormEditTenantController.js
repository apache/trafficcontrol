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
 * @param {*} tenant
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/TenantService")} tenantService
 * @param {import("../../../../models/MessageModel")} messageModel
 */
var FormEditTenantController = function(tenant, $scope, $controller, $uibModal, $anchorScroll, locationUtils, tenantService, messageModel) {

    // extends the FormTenantController to inherit common methods
    angular.extend(this, $controller('FormTenantController', { tenant: tenant, $scope: $scope }));

    var deleteTenant = function(tenant) {
        tenantService.deleteTenant(tenant.id)
            .then(function(result) {
                messageModel.setMessages(result.data.alerts, true);
                locationUtils.navigateToPath('/tenants');
            });
    };

    $scope.tenantName = angular.copy(tenant.name);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(tenant) {
        tenantService.updateTenant(tenant).
            then(function() {
                $scope.tenantName = angular.copy(tenant.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(tenant) {
        var params = {
            title: 'Delete Tenant: ' + tenant.name,
            key: tenant.name
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
            deleteTenant(tenant);
        }, function () {
            // do nothing
        });
    };

};

FormEditTenantController.$inject = ['tenant', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'tenantService', 'messageModel'];
module.exports = FormEditTenantController;
