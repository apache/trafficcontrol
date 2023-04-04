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
 * @param {*} cacheGroup
 * @param {*} types
 * @param {*} cacheGroups
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/CacheGroupService")} cacheGroupService
 * @param {import("../../../../models/MessageModel")} messageModel
 */
var FormEditCacheGroupController = function(cacheGroup, types, cacheGroups, $scope, $controller, $uibModal, $anchorScroll, locationUtils, cacheGroupService, messageModel) {

    // extends the FormCacheGroupController to inherit common methods
    angular.extend(this, $controller('FormCacheGroupController', { cacheGroup: cacheGroup, types: types, cacheGroups: cacheGroups, $scope: $scope }));

    $scope.cacheGroupName = angular.copy(cacheGroup.name);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    var deleteCacheGroup = function(cacheGroup) {
        cacheGroupService.deleteCacheGroup(cacheGroup.id)
            .then(function(result) {
                messageModel.setMessages(result.alerts, true);
                locationUtils.navigateToPath('/cache-groups');
            });
    };

    var queueServerUpdates = function(cacheGroup, cdnId) {
        cacheGroupService.queueServerUpdates(cacheGroup.id, cdnId);
    };

    var clearServerUpdates = function(cacheGroup, cdnId) {
        cacheGroupService.clearServerUpdates(cacheGroup.id, cdnId);
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

    $scope.confirmQueueServerUpdates = function(cacheGroup) {
        var params = {
            title: 'Queue Server Updates: ' + cacheGroup.name,
            message: "Please select a CDN"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
            controller: 'DialogSelectController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function(cdnService) {
                    return cdnService.getCDNs();
                }
            }
        });
        modalInstance.result.then(function(cdn) {
            queueServerUpdates(cacheGroup, cdn.id);
        }, function () {
            // do nothing
        });
    };

    $scope.confirmClearServerUpdates = function(cacheGroup) {
        var params = {
            title: 'Clear Server Updates: ' + cacheGroup.name,
            message: "Please select a CDN"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
            controller: 'DialogSelectController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function(cdnService) {
                    return cdnService.getCDNs();
                }
            }
        });
        modalInstance.result.then(function(cdn) {
            clearServerUpdates(cacheGroup, cdn.id);
        }, function () {
            // do nothing
        });
    };

};

FormEditCacheGroupController.$inject = ['cacheGroup', 'types', 'cacheGroups', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'cacheGroupService', 'messageModel'];
module.exports = FormEditCacheGroupController;
