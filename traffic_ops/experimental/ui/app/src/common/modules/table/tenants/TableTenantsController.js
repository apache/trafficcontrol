var TableTenantsController = function(tenants, $scope, locationUtils) {

    $scope.tenants = tenants;

    $scope.editTenant = function(id) {
        locationUtils.navigateToPath('/admin/tenants/' + id + '/edit');
    };

    $scope.createTenant = function() {
        locationUtils.navigateToPath('/admin/tenants/new');
    };

    angular.element(document).ready(function () {
        $('#tenantsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableTenantsController.$inject = ['tenants', '$scope', 'locationUtils'];
module.exports = TableTenantsController;