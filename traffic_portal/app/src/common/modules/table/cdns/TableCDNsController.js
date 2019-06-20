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

var TableCDNsController = function(cdns, $location, $scope, $state, $uibModal, $window, locationUtils, cdnService, messageModel) {

    var queueServerUpdates = function(cdn) {
        cdnService.queueServerUpdates(cdn.id);
    };

    var clearServerUpdates = function(cdn) {
        cdnService.clearServerUpdates(cdn.id);
    };

    var deleteCDN = function(cdn) {
        cdnService.deleteCDN(cdn.id)
            .then(function(result) {
                messageModel.setMessages(result.alerts, false);
                $scope.refresh();
            });
    };

    var confirmQueueServerUpdates = function(cdn) {
        var params = {
            title: 'Queue Server Updates: ' + cdn.name,
            message: 'Are you sure you want to queue server updates for all ' + cdn.name + ' servers?'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            queueServerUpdates(cdn);
        }, function () {
            // do nothing
        });
    };

    var confirmClearServerUpdates = function(cdn) {
        var params = {
            title: 'Clear Server Updates: ' + cdn.name,
            message: 'Are you sure you want to clear server updates for all ' + cdn.name + ' servers?'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            clearServerUpdates(cdn);
        }, function () {
            // do nothing
        });
    };

    var confirmDelete = function(cdn) {
        var params = {
            title: 'Delete CDN: ' + cdn.name,
            key: cdn.name
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
            deleteCDN(cdn);
        }, function () {
            // do nothing
        });
    };

    $scope.cdns = cdns;

    $scope.contextMenuItems = [
        {
            text: 'Open in New Tab',
            click: function ($itemScope) {
                $window.open('/#!/cdns/' + $itemScope.cdn.id, '_blank');
            }
        },
        null, // Dividier
        {
            text: 'Edit',
            click: function ($itemScope) {
                $scope.editCDN($itemScope.cdn.id);
            }
        },
        {
            text: 'Delete',
            click: function ($itemScope) {
                confirmDelete($itemScope.cdn);
            }
        },
        null, // Dividier
        {
            text: 'Diff Snapshot',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/cdns/' + $itemScope.cdn.id + '/config/changes');
            }
        },
        null, // Dividier
        {
            text: 'Queue Server Updates',
            click: function ($itemScope) {
                confirmQueueServerUpdates($itemScope.cdn);
            }
        },
        {
            text: 'Clear Server Updates',
            click: function ($itemScope) {
                confirmClearServerUpdates($itemScope.cdn);
            }
        },
        null, // Dividier
        {
            text: 'Manage DNSSEC Keys',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/cdns/' + $itemScope.cdn.id + '/dnssec-keys');
            }
        },
        {
            text: 'Manage Federations',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/cdns/' + $itemScope.cdn.id + '/federations');
            }
        },
        {
            text: 'Manage Delivery Services',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/cdns/' + $itemScope.cdn.id + '/delivery-services');
            }
        },
        {
            text: 'Manage Profiles',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/cdns/' + $itemScope.cdn.id + '/profiles');
            }
        },
        {
            text: 'Manage Servers',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/cdns/' + $itemScope.cdn.id + '/servers');
            }
        }
    ];

    $scope.editCDN = function(id) {
        locationUtils.navigateToPath('/cdns/' + id);
    };

    $scope.createCDN = function() {
        locationUtils.navigateToPath('/cdns/new');
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    angular.element(document).ready(function () {
        $('#cdnsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "aaSorting": []
        });
    });

};

TableCDNsController.$inject = ['cdns', '$location', '$scope', '$state', '$uibModal', '$window', 'locationUtils', 'cdnService', 'messageModel'];
module.exports = TableCDNsController;
