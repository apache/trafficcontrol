var TableDivisionsController = function(divisions, $scope, $location) {

    $scope.divisions = divisions;

    $scope.editDivision = function(id) {
        $location.path($location.path() + '/' + id);
    };

    angular.element(document).ready(function () {
        $('#divisionsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableDivisionsController.$inject = ['divisions', '$scope', '$location'];
module.exports = TableDivisionsController;