var TableServersController = function(servers, $scope, $location) {

    $scope.servers = servers;

    $scope.editServer = function(id) {
        $location.path($location.path() + id);
    };

    angular.element(document).ready(function () {
        $('#serversTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 100
        });
    });

};

TableServersController.$inject = ['servers', '$scope', '$location'];
module.exports = TableServersController;