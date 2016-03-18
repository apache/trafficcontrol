var TableLocationsController = function(locations, $scope, locationUtils) {

    $scope.locations = locations;

    $scope.editLocation = function(id) {
        locationUtils.navigateToPath('/admin/locations/' + id + '/edit')
    };

    $scope.createLocation = function() {
        locationUtils.navigateToPath('/admin/locations/new')
    };

    angular.element(document).ready(function () {
        $('#locationsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableLocationsController.$inject = ['locations', '$scope', 'locationUtils'];
module.exports = TableLocationsController;