var TableRegionsController = function(regions, $scope) {

    $scope.regions = regions;

    angular.element(document).ready(function () {
        $('#regionsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableRegionsController.$inject = ['regions', '$scope'];
module.exports = TableRegionsController;