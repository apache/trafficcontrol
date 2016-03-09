var TableServersController = function(servers, $scope) {

    $scope.servers = servers;

    angular.element(document).ready(function () {
        $('#serversTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 100
        });
    });

};

TableServersController.$inject = ['servers', '$scope'];
module.exports = TableServersController;