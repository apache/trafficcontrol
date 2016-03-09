var TableLocationsController = function(locations, $scope) {

    $scope.locations = locations;

    angular.element(document).ready(function () {
        $('#locationsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableLocationsController.$inject = ['locations', '$scope'];
module.exports = TableLocationsController;