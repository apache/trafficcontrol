var TableTenantsController = function(tenants, $scope, $location) {

    $scope.tenants = tenants;

    $scope.editTenant = function(id) {
        $location.path($location.path() + id);
    };

    angular.element(document).ready(function () {
        $('#tenantsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableTenantsController.$inject = ['tenants', '$scope', '$location'];
module.exports = TableTenantsController;