angular.module('app', [])

.controller('InitCtrl', function($scope, $http, $log) {
    $http.get('http://127.0.0.1:8080/request/').then(function(resp) {
        $scope.tables = resp.data;
        // For JSON responses, resp.data contains the result
    }, function(err) {
        console.error('ERR', err);
        // err.status will contain the status code
    })

    //GET
    $scope.get = function(table) {
        var tableName = angular.copy(table);

        $http.get('http://127.0.0.1:8080/request/' + tableName).then(function(resp) {
            $scope.columns = resp.data;
            // For JSON responses, resp.data contains the result
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })

        $http.get('http://127.0.0.1:8080/api/' + tableName).then(function(resp) {
            $scope.rows = resp.data.response;
            console.log(resp.data);
            if (resp.data.error != "") {
                alert(resp.data.error);
            }
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })
    }

    //GET
    $scope.update = function(table, parameters) {
        var tableName = angular.copy(table);

        if (typeof parameters !== 'undefined') {
            $http.get('http://127.0.0.1:8080/request/' + tableName + "?" + parameters).then(function(resp) {
                console.log(resp.data);
                $scope.columns = resp.data;
                // For JSON responses, resp.data contains the result
            }, function(err) {
                console.error('ERR', err);
                // err.status will contain the status code
            })

            $http.get('http://127.0.0.1:8080/api/' + tableName + "?" + parameters).then(function(resp) {
                $scope.rows = resp.data.response;
                if (resp.data.error != "") {
                    alert(resp.data.error);
                }

            }, function(err) {
                console.error('ERR', err);
                // err.status will contain the status code
            })
        } else {
            $scope.get(table);
        }
    }

    //DELETE
    $scope.delete = function(table, row) {
        $http.get('http://127.0.0.1:8080/request/' + table).then(function(resp) {
            console.log(resp.data);
            $scope.columns = resp.data;
            // For JSON responses, resp.data contains the result
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })

        $http.delete('http://127.0.0.1:8080/api/' + table + "/" + row.id).then(function(resp) {
            $scope.rows = resp.data.response;
            if (resp.data.error != "") {
                alert(resp.data.error);
            }

            //make table
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })
    }

    //POST QUERY
    $scope.postView = function(newView) {
        var viewArray = new Array(newView);

        //post it
        $http.post('http://127.0.0.1:8080/api/', viewArray).then(function(resp) {
            $scope.rows = resp.data.response;
            if (resp.data.error != "") {
                alert(resp.data.error);
            }

        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })

        //get columns
        $http.get('http://127.0.0.1:8080/request/' + newView.Name).then(function(resp) {
            console.log("COLUMNS: " + resp.data);
            $scope.columns = resp.data;
            // For JSON responses, resp.data contains the result
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })
    }

    $scope.post = function(table, row) {
        var rowArray = new Array(row);

        //post it
        console.log(table, row);
        $http.post('http://127.0.0.1:8080/api/' + table, rowArray).then(function(resp) {
            $scope.rows = resp.data.response;
            if (resp.data.error != "") {
                alert(resp.data.error);
            }

        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })

        //get columns
        $http.get('http://127.0.0.1:8080/request/' + table).then(function(resp) {
            $scope.columns = resp.data;
            // For JSON responses, resp.data contains the result
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })
    }

    //PUT
    $scope.put = function(table, row) {
        var rowArray = new Array(row);
        //post it
        $http.put('http://127.0.0.1:8080/api/' + table + "/" + row.id, rowArray).then(function(resp) {
            $scope.rows = resp.data.response;
            if (resp.data.error != "") {
                alert(resp.data.error);
            }

        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })

        //get columns
        $http.get('http://127.0.0.1:8080/request/' + table).then(function(resp) {
            $scope.columns = resp.data;
            // For JSON responses, resp.data contains the result
        }, function(err) {
            console.error('ERR', err);
            // err.status will contain the status code
        })
    }

})

