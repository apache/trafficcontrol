var TableLocationsController = function(locations, $scope, $location) {

    $scope.locations = locations;

    $scope.editLocation = function(id) {
        $location.path($location.path() + '/' + id);
    };

    angular.element(document).ready(function () {
        $('#locationsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableLocationsController.$inject = ['locations', '$scope', '$location'];
module.exports = TableLocationsController;