var TableDivisionsController = function(divisions, $scope, locationUtils) {

    $scope.divisions = divisions;

    $scope.editDivision = function(id) {
        locationUtils.navigateToPath('/admin/divisions/' + id + '/edit');
    };

    $scope.createDivision = function() {
        locationUtils.navigateToPath('/admin/divisions/new');
    };

    angular.element(document).ready(function () {
        $('#divisionsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableDivisionsController.$inject = ['divisions', '$scope', 'locationUtils'];
module.exports = TableDivisionsController;