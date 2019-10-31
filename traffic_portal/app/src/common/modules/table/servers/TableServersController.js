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

var TableServersController = function(servers, $scope, $state, $uibModal, $window, dateUtils, locationUtils, serverUtils, cdnService, serverService, statusService, propertiesModel, messageModel) {

    let serversTable;

    var getStatuses = function() {
        statusService.getStatuses()
            .then(function(result) {
                $scope.statuses = result;
            });
    };

    var queueServerUpdates = function(server) {
        serverService.queueServerUpdates(server.id)
            .then(
                function() {
                    $scope.refresh();
                }
            );
    };

    var clearServerUpdates = function(server) {
        serverService.clearServerUpdates(server.id)
            .then(
                function() {
                    $scope.refresh();
                }
            );
    };

    var queueCDNServerUpdates = function(cdnId) {
        cdnService.queueServerUpdates(cdnId)
            .then(
                function() {
                    $scope.refresh();
                }
            );
    };

    var clearCDNServerUpdates = function(cdnId) {
        cdnService.clearServerUpdates(cdnId)
            .then(
                function() {
                    $scope.refresh();
                }
            );
    };

    var confirmDelete = function(server) {
        var params = {
            title: 'Delete Server: ' + server.hostName,
            key: server.hostName
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
            deleteServer(server);
        }, function () {
            // do nothing
        });
    };

    var deleteServer = function(server) {
        serverService.deleteServer(server.id)
            .then(function(result) {
                messageModel.setMessages(result.alerts, false);
                $scope.refresh();
            });
    };

    var confirmStatusUpdate = function(server) {
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/status/dialog.select.status.tpl.html',
            controller: 'DialogSelectStatusController',
            size: 'md',
            resolve: {
                server: function() {
                    return server;
                },
                statuses: function() {
                    return $scope.statuses;
                }
            }
        });
        modalInstance.result.then(function(status) {
            updateStatus(status, server);
        }, function () {
            // do nothing
        });
    };

    var updateStatus = function(status, server) {
        serverService.updateStatus(server.id, { status: status.id, offlineReason: status.offlineReason })
            .then(
                function(result) {
                    messageModel.setMessages(result.data.alerts, false);
                    $scope.refresh();
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    $scope.servers = servers;

    $scope.columns = [
        { "name": "Cache Group", "visible": true, "searchable": true },
        { "name": "CDN", "visible": true, "searchable": true },
        { "name": "Domain", "visible": true, "searchable": true },
        { "name": "Host", "visible": true, "searchable": true },
        { "name": "HTTPS Port", "visible": false, "searchable": false },
        { "name": "ID", "visible": false, "searchable": false },
        { "name": "ILO IP Address", "visible": true, "searchable": true },
        { "name": "ILO IP Gateway", "visible": false, "searchable": false },
        { "name": "ILO IP Netmask", "visible": false, "searchable": false },
        { "name": "ILO Username", "visible": false, "searchable": false },
        { "name": "Interface Name", "visible": false, "searchable": false },
        { "name": "IPv6 Address", "visible": true, "searchable": true },
        { "name": "IPv6 Gateway", "visible": false, "searchable": false },
        { "name": "Last Updated", "visible": false, "searchable": false },
        { "name": "Mgmt IP Address", "visible": false, "searchable": false },
        { "name": "Mgmt IP Gateway", "visible": false, "searchable": false },
        { "name": "Mgmt IP Netmask", "visible": false, "searchable": false },
        { "name": "Network Gateway", "visible": false, "searchable": false },
        { "name": "Network IP", "visible": true, "searchable": true },
        { "name": "Network MTU", "visible": false, "searchable": false },
        { "name": "Network Subnet", "visible": false, "searchable": false },
        { "name": "Offline Reason", "visible": false, "searchable": false },
        { "name": "Phys Location", "visible": true, "searchable": true },
        { "name": "Profile", "visible": true, "searchable": true },
        { "name": "Rack", "visible": false, "searchable": false },
        { "name": "Reval Pending", "visible": false, "searchable": false },
        { "name": "Router Hostname", "visible": false, "searchable": false },
        { "name": "Router Port Name", "visible": false, "searchable": false },
        { "name": "Status", "visible": true, "searchable": true },
        { "name": "TCP Port", "visible": false, "searchable": false },
        { "name": "Type", "visible": true, "searchable": true },
        { "name": "Update Pending", "visible": true, "searchable": true }
    ];

    $scope.contextMenuItems = [
        {
            text: 'Open in New Tab',
            click: function ($itemScope) {
                $window.open('/#!/servers/' + $itemScope.s.id, '_blank');
            }
        },
        null, // Divider
        {
            text: 'Navigate to Server FQDN',
            click: function ($itemScope) {
                $window.open('http://' + $itemScope.s.hostName + '.' + $itemScope.s.domainName, '_blank');
            }
        },
        null, // Divider
        {
            text: 'Edit',
            click: function ($itemScope) {
                $scope.editServer($itemScope.s.id);
            }
        },
        {
            text: 'Delete',
            click: function ($itemScope) {
                confirmDelete($itemScope.s);
            }
        },
        null, // Divider
        {
            text: 'Update Status',
            click: function ($itemScope) {
                confirmStatusUpdate($itemScope.s);
            }
        },
        {
            text: 'Queue Server Updates',
            displayed: function ($itemScope) {
                return serverUtils.isCache($itemScope.s) && !$itemScope.s.updPending;
            },
            click: function ($itemScope) {
                queueServerUpdates($itemScope.s);
            }
        },
        {
            text: 'Clear Server Updates',
            displayed: function ($itemScope) {
                return serverUtils.isCache($itemScope.s) && $itemScope.s.updPending;
            },
            click: function ($itemScope) {
                clearServerUpdates($itemScope.s);
            }
        },
        {
            text: 'Show Charts',
            displayed: function () {
                return propertiesModel.properties.servers.charts.show;
            },
            hasBottomDivider: function () {
                return true;
            },
            hasTopDivider: function () {
                return true;
            },
            click: function ($itemScope) {
                $window.open(propertiesModel.properties.servers.charts.baseUrl + $itemScope.s.hostName, '_blank');
            }
        },
        {
            text: 'Manage Capabilities',
            displayed: function ($itemScope) {
                return serverUtils.isCache($itemScope.s);
            },
            hasTopDivider: function () {
                return true;
            },
            click: function ($itemScope) {
                locationUtils.navigateToPath('/servers/' + $itemScope.s.id + '/capabilities');
            }
        },
        {
            text: 'Manage Delivery Services',
            displayed: function ($itemScope) {
                return serverUtils.isEdge($itemScope.s);
            },
            click: function ($itemScope) {
                locationUtils.navigateToPath('/servers/' + $itemScope.s.id + '/delivery-services');
            }
        },
        {
            text: 'View Config Files',
            displayed: function ($itemScope) {
                return serverUtils.isCache($itemScope.s);
            },
            click: function ($itemScope) {
                locationUtils.navigateToPath('/servers/' + $itemScope.s.id + '/config-files');
            }
        }
    ];

    $scope.editServer = function(id) {
        locationUtils.navigateToPath('/servers/' + id);
    };

    $scope.createServer = function() {
        locationUtils.navigateToPath('/servers/new');
    };

    $scope.confirmCDNQueueServerUpdates = function(cdn) {
        var params;
        if (cdn) {
            params = {
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
                queueCDNServerUpdates(cdn.id);
            }, function () {
                // do nothing
            });
        } else {
            params = {
                title: 'Queue Server Updates',
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
                queueCDNServerUpdates(cdn.id);
            }, function () {
                // do nothing
            });
        }
    };

    $scope.confirmCDNClearServerUpdates = function(cdn) {
        var params;
        if (cdn) {
            params = {
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
                clearCDNServerUpdates(cdn.id);
            }, function () {
                // do nothing
            });


        } else {
            params = {
                title: 'Clear Server Updates',
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
                clearCDNServerUpdates(cdn.id);
            }, function () {
                // do nothing
            });
        }
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    $scope.toggleVisibility = function(colName) {
        const col = serversTable.column(colName + ':name');
        col.visible(!col.visible());
        serversTable.rows().invalidate().draw();
    };

    $scope.ssh = serverUtils.ssh;

    $scope.isOffline = serverUtils.isOffline;

    $scope.offlineReason = serverUtils.offlineReason;

    $scope.getRelativeTime = dateUtils.getRelativeTime;

    $scope.navigateToPath = locationUtils.navigateToPath;

    var init = function () {
        getStatuses();
    };
    init();

    angular.element(document).ready(function () {
        serversTable = $('#serversTable').DataTable({
            "lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "aaSorting": [],
            "columns": $scope.columns,
            "colReorder": {
                realtime: false
            },
            "initComplete": function(settings, json) {
                try {
                    // need to create the show/hide column checkboxes and bind to the current visibility
                    $scope.columns = JSON.parse(localStorage.getItem('DataTables_serversTable_/')).columns;
                } catch (e) {
                    console.error("Failure to retrieve required column info from localStorage (key=DataTables_serversTable_/):", e);
                }
            }
        });
    });

};

TableServersController.$inject = ['servers', '$scope', '$state', '$uibModal', '$window', 'dateUtils', 'locationUtils', 'serverUtils', 'cdnService', 'serverService', 'statusService', 'propertiesModel', 'messageModel'];
module.exports = TableServersController;
