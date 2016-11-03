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

var FormEditCacheGroupController = function(cacheGroup, $scope, $controller, $uibModal, $anchorScroll, locationUtils, cacheGroupService) {

    // extends the FormCacheGroupController to inherit common methods
    angular.extend(this, $controller('FormCacheGroupController', { cacheGroup: cacheGroup, $scope: $scope }));

    var deleteCacheGroup = function(cacheGroup) {
        cacheGroupService.deleteCacheGroup(cacheGroup.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/cache-groups');
            });
    };

    $scope.cacheGroupName = angular.copy(cacheGroup.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(cacheGroup) {
        cacheGroupService.updateCacheGroup(cacheGroup).
            then(function() {
                $scope.cacheGroupName = angular.copy(cacheGroup.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(cacheGroup) {
        var params = {
            title: 'Delete Cache Group: ' + cacheGroup.name,
            key: cacheGroup.name
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
            deleteCacheGroup(cacheGroup);
        }, function () {
            // do nothing
        });
    };

};

FormEditCacheGroupController.$inject = ['cacheGroup', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'cacheGroupService'];
module.exports = FormEditCacheGroupController;
