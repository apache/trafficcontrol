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

var TableAssignServerSCsController = function(server, serverCapabilities, assignedSCs, $scope, $uibModalInstance) {

    $scope.selectedSCs = [];

    $scope.server = server;

    /** @type CGC.ColumnDefinition */
    $scope.columns = [
        {
            headerName: "Server Capability",
            field: "name",
            checkboxSelection: true,
            headerCheckboxSelection: true,
        }
    ];

    $scope.serverCapabilities = serverCapabilities.map(SCs => {
        let isAssigned = assignedSCs.find(assignedSCs => assignedSCs.serverCapability === SCs.name);
        if (isAssigned) {
            SCs['selected'] = true;
        }
        return SCs;
    });

    $scope.submit = function() {
        const selectedSCNames = this.selectedSCs.map(sc => sc["name"]);
        $uibModalInstance.close(selectedSCNames);
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

TableAssignServerSCsController.$inject = ['server', 'serverCapabilities', 'assignedSCs', '$scope', '$uibModalInstance'];
module.exports = TableAssignServerSCsController;

