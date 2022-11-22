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
 * @param {*} cacheGroups
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IWindowService} $window
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/CacheGroupService")} cacheGroupService
 * @param {import("../../../models/MessageModel")} messageModel
 */
var TableCacheGroupsController = function(cacheGroups, $scope, $state, $uibModal, $window, locationUtils, cacheGroupService, messageModel) {

    let cacheGroupsTable;

    var queueServerUpdates = function(cacheGroup, cdnId) {
        cacheGroupService.queueServerUpdates(cacheGroup.id, cdnId);
    };

    var clearServerUpdates = function(cacheGroup, cdnId) {
        cacheGroupService.clearServerUpdates(cacheGroup.id, cdnId);
    };

    var deleteCacheGroup = function(cacheGroup) {
        cacheGroupService.deleteCacheGroup(cacheGroup.id)
            .then(function(result) {
                messageModel.setMessages(result.alerts, false);
                $scope.refresh();
            });
    };

    var confirmQueueServerUpdates = function(cacheGroup) {
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

    var confirmClearServerUpdates = function(cacheGroup) {
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

    var confirmDelete = function(cacheGroup) {
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


    $scope.cacheGroups = cacheGroups;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.columns = [
        { "name": "Name", "visible": true, "searchable": true },
        { "name": "Short Name", "visible": true, "searchable": true },
        { "name": "Type", "visible": true, "searchable": true },
        { "name": "1st Parent", "visible": true, "searchable": true },
        { "name": "2nd Parent", "visible": true, "searchable": true },
        { "name": "Latitude", "visible": true, "searchable": true },
        { "name": "Longitude", "visible": true, "searchable": true }
    ];

    $scope.contextMenuItems = [
        {
            text: 'Open in New Tab',
            click: function ($itemScope) {
                $window.open('/#!/cache-groups/' + $itemScope.cg.id, '_blank');
            }
        },
        null, // Dividier
        {
            text: 'Edit',
            click: function ($itemScope) {
                $scope.editCacheGroup($itemScope.cg.id);
            }
        },
        {
            text: 'Delete',
            click: function ($itemScope) {
                confirmDelete($itemScope.cg);
            }
        },
        null, // Dividier
        {
            text: 'Queue Server Updates',
            click: function ($itemScope) {
                confirmQueueServerUpdates($itemScope.cg);
            }
        },
        {
            text: 'Clear Server Updates',
            click: function ($itemScope) {
                confirmClearServerUpdates($itemScope.cg);
            }
        },
        null, // Dividier
        {
            text: 'Manage ASNs',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/cache-groups/' + $itemScope.cg.id + '/asns');
            }
        },
        {
            text: 'Manage Servers',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/cache-groups/' + $itemScope.cg.id + '/servers');
            }
        }
    ];

    $scope.editCacheGroup = function(id) {
        locationUtils.navigateToPath('/cache-groups/' + id);
    };

    $scope.createCacheGroup = function() {
        locationUtils.navigateToPath('/cache-groups/new');
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    $scope.toggleVisibility = function(colName) {
        const col = cacheGroupsTable.column(colName + ':name');
        col.visible(!col.visible());
        cacheGroupsTable.rows().invalidate().draw();
    };

    angular.element(document).ready(function () {
        cacheGroupsTable = $('#cacheGroupsTable').DataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "aaSorting": [],
            "columns": $scope.columns,
            "initComplete": function(settings, json) {
                try {
                    // need to create the show/hide column checkboxes and bind to the current visibility
                    $scope.columns = JSON.parse(localStorage.getItem('DataTables_cacheGroupsTable_/')).columns;
                } catch (e) {
                    console.error("Failure to retrieve required column info from localStorage (key=DataTables_cacheGroupsTable_/):", e);
                }
            }
        });
    });

};

TableCacheGroupsController.$inject = ['cacheGroups', '$scope', '$state', '$uibModal', '$window', 'locationUtils', 'cacheGroupService', 'messageModel'];
module.exports = TableCacheGroupsController;
