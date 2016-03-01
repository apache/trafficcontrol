var TableTenantsController = function(tenants, $scope) {

    $scope.tenants = tenants;

    angular.element(document).ready(function () {
        $('#tenantsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableTenantsController.$inject = ['tenants', '$scope'];
module.exports = TableTenantsController;