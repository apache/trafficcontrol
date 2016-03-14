var TableRegionsController = function(regions, $scope, $location) {

    $scope.regions = regions;

    $scope.editRegion = function(id) {
        $location.path($location.path() + '/' + id);
    };

    angular.element(document).ready(function () {
        $('#regionsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableRegionsController.$inject = ['regions', '$scope', '$location'];
module.exports = TableRegionsController;