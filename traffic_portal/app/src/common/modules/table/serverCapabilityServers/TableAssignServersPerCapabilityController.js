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

/** @typedef { import('../../../../common/modules/table/agGrid/CommonGridController').CGC } CGC */

var TableAssignServersPerCapabilityController = function(servers, serverCapability, assignedServers, $scope, $uibModalInstance) {

    $scope.selectedServers = [];

    $scope.serverCapability = serverCapability;

    /** @type CGC.ColumnDefinition */
    $scope.columns = [
        {
            headerName: "Servers",
            field: "hostName",
            checkboxSelection: true,
            headerCheckboxSelection: true,
            sort: "asc",
        },
        {
            headerName: "Type",
            field: "type",
            hide: false
        },
        {
            headerName: "CDN",
            field: "cdn",
            hide: false
        }
    ];

    $scope.servers = servers.map(server => {
        let isAssigned = assignedServers.find(assignedServers => assignedServers.serverId === server.id);
        if (isAssigned) {
            server['selected'] = true;
        }
        return server;
    });

    $scope.submit = function() {
        const selectedServerIds = this.selectedServers.map(mspc => mspc["id"]);
        $uibModalInstance.close(selectedServerIds);
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };

    /** @type CGC.GridSettings */
    $scope.gridOptions = {
        selectRows: true,
        selectionProperty: "selected"
    };
};

TableAssignServersPerCapabilityController.$inject = ['servers', 'serverCapability', 'assignedServers', '$scope', '$uibModalInstance'];
module.exports = TableAssignServersPerCapabilityController;
