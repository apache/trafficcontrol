var TableDivisionsController = function(divisions, $scope) {

    $scope.divisions = divisions.response;

    angular.element(document).ready(function () {
        $('#divisionsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableDivisionsController.$inject = ['divisions', '$scope'];
module.exports = TableDivisionsController;