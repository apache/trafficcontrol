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


var TableServersController = function(servers, $scope, $state, $uibModal, $window, dateUtils, locationUtils, serverUtils, cdnService, serverService, statusService, propertiesModel, messageModel, userModel) {

    // browserify can't handle classes...
    function SSHCellRenderer() {}
    SSHCellRenderer.prototype.init = function(params) {
        this.eGui = document.createElement("A");
        this.eGui.href = "ssh://" + userModel.user.username + "@" + params.value;
        this.eGui.setAttribute("target", "_blank");
        this.eGui.textContent = params.value;
    };
    SSHCellRenderer.prototype.getGui = function() {return this.eGui;};

    function UpdateCellRenderer() {}
    UpdateCellRenderer.prototype.init = function(params) {
        this.eGui = document.createElement("I");
        this.eGui.setAttribute("aria-hidden", "true");
        this.eGui.setAttribute("title", String(params.value));
        this.eGui.classList.add("fa", "fa-lg");
        if (params.value) {
            this.eGui.classList.add("fa-check");
        } else {
            this.eGui.classList.add("fa-clock-o");
        }
    }
    UpdateCellRenderer.prototype.getGui = function() {return this.eGui;};

    function offlineReasonTooltip(params) {
        if (!params.value || !serverUtils.isOffline(params.value)) {
            return;
        }
        return params.data.offlineReason;
    }

    function dateCellFormatter(params) {
        return dateUtils.getRelativeTime(params.value);
    }

    let serversTable;

    function editServer(params) {
        $scope.editServer(params.data.id);
        // Event is outside the digest cycle, so we need this to trigger one.
        $scope.$apply();
    }

    const agColumns = [
        {
            headerName: "Cache Group",
            field: "cachegroup",
            hide: false,
            searchable: true
        },
        {
            headerName: "CDN",
            field: "cdn",
            hide: false,
            searchable: true
        },
        {
            headerName: "Domain",
            field: "domainName",
            hide: false,
            searchable: true
        },
        {
            headerName: "Host",
            field: "hostName",
            hide: false,
            searchable: true
        },
        {
            headerName: "HTTPS Port",
            field: "httpsPort",
            hide: true,
            searchable: false,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "ID",
            field: "id",
            hide: true,
            searchable: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "ILO IP Address",
            field: "iloIpAddress",
            hide: true,
            searchable: true,
            cellRenderer: "sshCellRenderer",
            onCellClicked: null
        },
        {
            headerName: "ILO IP Gateway",
            field: "iloIpGateway",
            hide: true,
            searchable: false,
            cellRenderer: "sshCellRenderer",
            onCellClicked: null
        },
        {
            headerName: "ILO IP Netmask",
            field: "iloIpNetmask",
            hide: true,
            searchable: false
        },
        {
            headerName: "ILO Username",
            field: "iloUsername",
            hide: true,
            searchable: false
        },
        {
            headerName: "Interface Name",
            field: "interfaceName",
            hide: true,
            searchable: false
        },
        {
            headerName: "IPv6 Address",
            field: "ipv6Address",
            hide: false,
            searchable: true
        },
        {
            headerName: "IPv6 Gateway",
            field: "ipv6Gateway",
            hide: true,
            searchable: false
        },
        {
            headerName: "Last Updated",
            field: "lastUpdated",
            hide: true,
            searchable: false,
            filter: "agDateColumnFilter",
            valueFormatter: dateCellFormatter
        },
        {
            headerName: "Mgmt IP Address",
            field: "mgmtIpAddress",
            hide: true,
            searchable: false
        },
        {
            headerName: "Mgmt IP Gateway",
            field: "mgmtIpGateway",
            hide: true,
            searchable: false,
            filter: true,
            cellRenderer: "sshCellRenderer",
            onCellClicked: null
        },
        {
            headerName: "Mgmt IP Netmask",
            field: "mgmtIpNetmask",
            hide: true,
            searchable: false,
            filter: true,
            cellRenderer: "sshCellRenderer",
            onCellClicked: null
        },
        {
            headerName: "Network Gateway",
            field: "ipGateway",
            hide: true,
            searchable: true,
            filter: true,
            cellRenderer: "sshCellRenderer",
            onCellClicked: null
        },
        {
            headerName: "Network IP",
            field: "ipAddress",
            hide: false,
            searchable: true,
            filter: true,
            cellRenderer: "sshCellRenderer",
            onCellClicked: null
        },
        {
            headerName: "Network MTU",
            field: "interfaceMtu",
            hide: true,
            searchable: false,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Network Subnet",
            field: "ipNetmask",
            hide: true,
            searchable: false
        },
        {
            headerName: "Offline Reason",
            field: "offlineReason",
            hide: true,
            searchable: false
        },
        {
            headerName: "Phys Location",
            field: "physLocation",
            hide: true,
            searchable: true
        },
        {
            headerName: "Profile",
            field: "profile",
            hide: false,
            searchable: true
        },
        {
            headerName: "Rack",
            field: "rack",
            hide: true,
            searchable: false
        },
        {
            headerName: "Reval Pending",
            field: "revalPending",
            hide: true,
            searchable: false,
            filter: true,
            cellRenderer: "updateCellRenderer"
        },
        {
            headerName: "Router Hostname",
            field: "routerHostName",
            hide: true,
            searchable: false
        },
        {
            headerName: "Router Port Name",
            field: "routerPortName",
            hide: true,
            searchable: false
        },
        {
            headerName: "Status",
            field: "status",
            hide: false,
            searchable: true,
            tooltip: offlineReasonTooltip
        },
        {
            headerName: "TCP Port",
            field: "tcpPort",
            hide: true,
            searchable: false
        },
        {
            headerName: "Type",
            field: "type",
            hide: false,
            searchable: true
        },
        {
            headerName: "Update Pending",
            field: "updPending",
            hide: false,
            searchable: true,
            filter: true,
            cellRenderer: "updateCellRenderer"
        }
    ];

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

    console.time("server map");
    $scope.servers = servers.map(function(x){x.lastUpdated = x.lastUpdated ? new Date(x.lastUpdated.replace("+00", "Z")) : x.lastUpdated;});
    console.timeEnd("server map");

    $scope.columns = [];

    $scope.gridOptions = {
        components: {
            sshCellRenderer: SSHCellRenderer,
            updateCellRenderer: UpdateCellRenderer
        },
        columnDefs: agColumns,
        defaultColDef: {
            filter: true,
            onCellClicked: editServer,
            sortable: true,
            resizable: true
        },
        rowData: servers,
        pagination: true,
        rowBuffer: 0,
        onGridReady: function(params) {
            params.api.sizeColumnsToFit();
        },
        menuTabs: 'columnsMenuTab',
        tooltipShowDelay: 500
    };

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
                return serverUtils.isEdge($itemScope.s) || serverUtils.isOrigin($itemScope.s);
            },
            hasTopDivider: function ($itemScope) {
                return !serverUtils.isCache($itemScope.s);
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

    $scope.toggleVisibility = function(col) {
        const visible = $scope.gridOptions.columnApi.getColumn(col).isVisible();
        $scope.gridOptions.columnApi.setColumnVisible(col, !visible);
        try {
            colsVisible = $scope.gridOptions.columnApi.getAllColumns().map(function(x){return [x.colId, x.isVisible()]});
            localStorage.setItem("servers_table_columns", JSON.stringify(colsVisible));
        } catch (e) {
            console.error("Failed to store column defs to local storage:", e);
        }
        $scope.gridOptions.api.sizeColumnsToFit();
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
        try {
            // need to create the show/hide column checkboxes and bind to the current visibility
            // TODO: figure out how to do this with getColumnState and setColumnState.
            const colstates = JSON.parse(localStorage.getItem("servers_table_columns"));
            if (colstates) {
                for (let i = 0; i<colstates.length; ++i) {
                    const colId = colstates[i][0];
                    const isVisible = colstates[i][1];
                    $scope.gridOptions.columnApi.setColumnVisible(colId, isVisible);
                }
            }
        } catch (e) {
            console.error("Failure to retrieve required column info from localStorage (key=servers_table_columns):", e);
        }
    });

};

TableServersController.$inject = ['servers', '$scope', '$state', '$uibModal', '$window', 'dateUtils', 'locationUtils', 'serverUtils', 'cdnService', 'serverService', 'statusService', 'propertiesModel', 'messageModel', "userModel"];
module.exports = TableServersController;
