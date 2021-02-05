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

var TableTenantsController = function(currentUserTenant, tenants, $scope, $state, $timeout, $uibModal, locationUtils, fileUtils, tenantUtils, tenantService, messageModel) {

    $scope.tenantTree = [];

    $scope.hasChildren = function(node) {
        return node.children.length > 0;
    };

    $scope.toggle = function(scope) {
        scope.toggle();
    };

    $scope.createTenant = function(parentId) {
        if (parentId) {
            locationUtils.navigateToPath('/tenants/new?parentId=' + parentId);
        } else {
            locationUtils.navigateToPath('/tenants/new');
        }
    };

    $scope.exportCSV = function() {
        fileUtils.convertToCSV(tenants, 'Tenants', ['id', 'lastUpdated', 'name', 'active', 'parentId', 'parentName']);
    };

    $scope.confirmDelete = function(tenant) {
        const params = {
            title: 'Delete Tenant: ' + tenant.name,
            key: tenant.name
        };
        const modalInstance = $uibModal.open({
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
            tenantService.deleteTenant(tenant.id)
                .then(function(result) {
                    messageModel.setMessages(result.data.alerts, false);
                    $state.reload();
                });
        }, function () {
            // do nothing
        });
    };

    let init = function() {
        $scope.tenants = tenantUtils.hierarchySort(tenantUtils.groupTenantsByParent(tenants), currentUserTenant.parentId, []);
        $scope.tenantTree = tenantUtils.convertToHierarchy($scope.tenants);
    };
    init();

};

TableTenantsController.$inject = ['currentUserTenant', 'tenants', '$scope', '$state', '$timeout', '$uibModal', 'locationUtils', 'fileUtils', 'tenantUtils', 'tenantService', 'messageModel'];
module.exports = TableTenantsController;
