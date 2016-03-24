var TableCDNsController = function(cdns, $scope, locationUtils) {

    $scope.cdns = cdns;

    $scope.editCDN = function(id) {
        locationUtils.navigateToPath('/admin/cdns/' + id + '/edit');
    };

    $scope.createCDN = function() {
        locationUtils.navigateToPath('/admin/cdns/new');
    };

    angular.element(document).ready(function () {
        $('#cdnsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableCDNsController.$inject = ['cdns', '$scope', 'locationUtils'];
module.exports = TableCDNsController;