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

angular.module('app', ['ngReactGrid'])

.controller('InitCtrl', function($scope, $http, $log, ngReactGridCheckbox) {
    var ipAddress = "localhost";
    //initialization
    $scope.colFilter = function(columns) {
		newColumns = {};
		for (var key in columns) {
			if (key != 'id' && key != 'last_updated') {
		newColumns[key] = columns[key];
}
}
return newColumns;
    }
    $scope.grid = {
        data: [],
        columnDefs: []
    }
    $scope.selections = [];
    getTableList();

    var checkboxGrid = new ngReactGridCheckbox($scope.selections, {
        batchToggle: true
    });

    $scope.filterRow = function(row) {
        for (var alias in row) {
            //is an alias
            if (!$scope.columns.hasOwnProperty(alias)) {
                //get row's antialiased column
                var column = $scope.aliasToColumnMap[alias];
                //use fkmap to get the value
				var curValue = row[alias];
				var curCol = $scope.columns[column];
if (curCol != undefined) {
				var fkValues = curCol["foreignKeyValues"];
				var value = fkValues[curValue];
                //append 
                row[column] = value;
                delete(row[alias]);
}
            }
        }
    }

    //get list of tables
    function getTableList() {
        $http.get('http://' + ipAddress + ':8080/request/').then(function(resp) {
            $scope.tables = resp.data;
        }, function(err) {
            console.error('ERR', err);
        })
    }

    function setTable(data) {
        $scope.newRow = {};

        if (data.error != "") {
            alert(data.error);
        }

        checkboxGrid.setVisibleCheckboxState(false);
        $scope.editEnabled = false;

        //set grid
        $scope.grid = {
            data: data.response,
            columnDefs: data.colWrappers.concat(checkboxGrid),
            horizontalScroll: true
        }

        $scope.isTable = data.isTable;
        $scope.columns = data.columns;
        $scope.aliasToColumnMap = {};
        for (var column in data.columns) {
            $scope.aliasToColumnMap[data.columns[column].colAlias] = column;
        }
    }

    $scope.clearCheckboxes = function() {
        checkboxGrid.setVisibleCheckboxState(false);
    }

    $scope.getColumnFromAlias = function(columnName) {
            for (var i = 0; i < $scope.columns.length; i++) {
                if ($scope.columns[i] == columnName) {
                    return $scope.columns[i];
                }
            }
        }
        //GET
    $scope.get = function(table) {
        $http.get('http://' + ipAddress + ':8080/api/' + table).then(function(resp) {
            setTable(resp.data);
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })
    }

    //GET
    $scope.update = function(table, parameters) {
        var tableName = angular.copy(table);

        if (typeof parameters !== 'undefined') {
            $http.get('http://' + ipAddress + ':8080/api/' + tableName + "?" + parameters).then(function(resp) {
                setTable(resp.data);
            }, function(err) {
                console.error('ERR', err);
            })
        } else {
            $scope.get(table);
        }
    }

    //DELETE
    $scope.delete = function(table, rows) {
        for (var i = 0; i < rows.length; i++) {
            $http.delete('http://' + ipAddress + ':8080/api/' + table + "/" + rows[i].id).then(function(resp) {
                setTable(resp.data);
            }, function(err) {
                console.error('ERR', err);
            })
        }
    }

    //DELETE
    $scope.deleteView = function(table) {
        $http.delete('http://' + ipAddress + ':8080/api/' + table).then(function(resp) {
            if (resp.data.error != "") {
                alert(resp.data.error);
            }

            location.reload();
            //make table
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })

        getTableList();
    }


    //POST QUERY
    $scope.postView = function(newView) {
        var viewArray = new Array(newView);

        $http.post('http://' + ipAddress + ':8080/api/', viewArray).then(function(resp) {
            if (resp.data.error != "") {
                alert(resp.data.error);
            }
            location.reload();
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })
    }

    $scope.post = function(table, row) {
        var rowArray = new Array(row);

        $http.post('http://' + ipAddress + ':8080/api/' + table, rowArray).then(function(resp) {
            setTable(resp.data);
        }, function(err) {
            console.error('ERR', err);
        })
    }

    //PUT
    $scope.put = function(table, rowArray) {
        //filter
        $http.put('http://' + ipAddress + ':8080/api/' + table, rowArray).then(function(resp) {
            setTable(resp.data);
        }, function(err) {
            console.error('ERR', err);
        })
    }
})
