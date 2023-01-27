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
 * @param {*} types
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/DateUtils")} dateUtils
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 */
var TableTypesController = function(types, $scope, $state, dateUtils, locationUtils) {

    let table;

    $scope.types = types;

    $scope.getRelativeTime = dateUtils.getRelativeTime;

    $scope.editType = function(id) {
        locationUtils.navigateToPath('/types/' + id);
    };

    $scope.createType = function() {
        locationUtils.navigateToPath('/types/new');
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    $scope.toggleVisibility = function(colName) {
        const col = table.column(colName + ':name');
        col.visible(!col.visible());
        table.rows().invalidate().draw();
    };

    angular.element(document).ready(function () {
        table = $('#typesTable').DataTable({
                "lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
                "iDisplayLength": 25,
                "aaSorting": [],
                "columns": [
                    { "name": "name", "visible": true, "searchable": true },
                    { "name": "description", "visible": true, "searchable": true },
                    { "name": "useInTable", "visible": true, "searchable": true },
                    { "name": "lastUpdated", "visible": false, "searchable": false }
                ],
                "initComplete": function(settings, json) {
                    try {
                        // need to create the show/hide column checkboxes and bind to the current visibility
                        $scope.columns = JSON.parse(localStorage.getItem('DataTables_typesTable_/')).columns;
                    } catch (e) {
                        console.error("Failure to retrieve required column info from localStorage (key=DataTables_typesTable_/):", e);
                    }
                }
            });
    });

};

TableTypesController.$inject = ['types', '$scope', '$state', 'dateUtils', 'locationUtils'];
module.exports = TableTypesController;
