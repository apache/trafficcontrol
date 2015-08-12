angular.module('app', ['ngReactGrid'])

.controller('InitCtrl', function($scope, $http, $log, ngReactGridCheckbox) {
    //initialization
    $scope.grid = {
        data: [],
        columnDefs: []
    }
    $scope.selections = [];
    getTableList();

    var checkboxGrid = new ngReactGridCheckbox($scope.selections, {
        batchToggle: true
    });
    //get list of tables
    function getTableList() {
        $http.get('http://127.0.0.1:8080/request/').then(function(resp) {
            $scope.tables = resp.data;
        }, function(err) {
            console.error('ERR', err);
        })
    }

    function setTable(data) {
        if (data.error != "") {
            alert(data.error);
        }

        checkboxGrid.setVisibleCheckboxState(false);

        //set grid
        $scope.grid = {
            data: data.response,
            columnDefs: data.columns.concat(checkboxGrid),
			height: 2000
        }
        $scope.isTable = data.isTable;

        //set columns for add row, edit row, etc
        var columns = [];
        for (var i = 0; i < data.columns.length; i++) {
            columns.push(data.columns[i].field);
        }

        $scope.columns = columns;

        //$scope.clearCheckboxes();
    }

    $scope.clearCheckboxes = function() {
        checkboxGrid.setVisibleCheckboxState(false);
    }

    //GET
    $scope.get = function(table) {
        $http.get('http://127.0.0.1:8080/api/' + table).then(function(resp) {
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
            $http.get('http://127.0.0.1:8080/api/' + tableName + "?" + parameters).then(function(resp) {
                setTable(resp.data);
            }, function(err) {
                console.error('ERR', err);
            })
        } else {
            $scope.get(table);
        }
    }

    //DELETE
    $scope.delete = function(table, row) {
        $http.delete('http://127.0.0.1:8080/api/' + table + "/" + row.id).then(function(resp) {
            setTable(resp.data);
        }, function(err) {
            console.error('ERR', err);
        })
    }

    //DELETE
    $scope.deleteView = function(table) {
        $http.delete('http://127.0.0.1:8080/api/' + table).then(function(resp) {
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

        $http.post('http://127.0.0.1:8080/api/', viewArray).then(function(resp) {
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

        $http.post('http://127.0.0.1:8080/api/' + table, rowArray).then(function(resp) {
            setTable(resp.data);
        }, function(err) {
            console.error('ERR', err);
        })
    }

    //PUT
    $scope.put = function(table, row) {
        var rowArray = new Array(row);
        $http.put('http://127.0.0.1:8080/api/' + table + "/" + row.id, rowArray).then(function(resp) {
            setTable(resp.data);
        }, function(err) {
            console.error('ERR', err);
        })
        console.log("hello");
    }
})
