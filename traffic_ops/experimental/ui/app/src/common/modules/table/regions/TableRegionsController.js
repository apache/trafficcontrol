var TableRegionsController = function(regions, $scope, locationUtils) {

    $scope.regions = regions;

    $scope.editRegion = function(id) {
        locationUtils.navigateToPath('/admin/regions/' + id + '/edit');
    };

    $scope.createRegion = function() {
        locationUtils.navigateToPath('/admin/regions/new');
    };

    angular.element(document).ready(function () {
        $('#regionsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableRegionsController.$inject = ['regions', '$scope', 'locationUtils'];
module.exports = TableRegionsController;