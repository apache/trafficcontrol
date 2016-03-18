var TableServersController = function(servers, $scope, locationUtils) {

    $scope.servers = servers;

    $scope.editServer = function(id) {
        locationUtils.navigateToPath('/configure/servers/' + id + '/edit');
    };

    $scope.createServer = function() {
        locationUtils.navigateToPath('/configure/servers/new');
    };

    angular.element(document).ready(function () {
        $('#serversTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 100
        });
    });

};

TableServersController.$inject = ['servers', '$scope', 'locationUtils'];
module.exports = TableServersController;