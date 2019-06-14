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

var TableTypesController = function(types, $scope, $state, $window, dateUtils, locationUtils) {

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

    angular.element(document).ready(function () {
        const table = $('#typesTable').DataTable({
                "lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
                "iDisplayLength": 25,
                "aaSorting": [],
                "columns": [
                    { "name": "name", "visible": true, "searchable": true },
                    { "name": "useInTable", "visible": true, "searchable": true },
                    { "name": "description", "visible": true, "searchable": true },
                    { "name": "lastUpdated", "visible": false, "searchable": false }
                ],
                "colReorder": {
                    realtime: false
                },
                "initComplete": function(settings, json) {
                    // need to bind the show/hide column checkboxes to the saved visibility
                    $scope.columns = JSON.parse($window.localStorage['DataTables_typesTable_/'])['columns'];
                    // also, need to reset column searchable to the column's saved visibility
                    $scope.columns.forEach(function(column, index) {
                        settings.aoColumns[index].bSearchable = column.visible;
                    });
                    // redraw so each column's searchable value is taken into account
                    this.api().rows().invalidate().draw();
                }
            });

        $('.column-settings input:checkbox').click(function() {
            const column = table.column($(this).data('column') + ':name');
            // toggle visibility of the selected table column
            column.visible(!column.visible());
            // hack alert: there is no api to set searchable on a column but if the column is visible, then it's searchable
            table.context[0].aoColumns[column.index()].bSearchable = column.visible();
            // redraw so the column's searchable value is taken into account
            table.rows().invalidate().draw();
        });
    });

};

TableTypesController.$inject = ['types', '$scope', '$state', '$window', 'dateUtils', 'locationUtils'];
module.exports = TableTypesController;
