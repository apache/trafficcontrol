var TableASNsController = function(asns, $scope, locationUtils) {

    $scope.asns = asns;

    $scope.editASN = function(id) {
        locationUtils.navigateToPath('/admin/asns/' + id + '/edit');
    };

    $scope.createASN = function() {
        locationUtils.navigateToPath('/admin/asns/new');
    };

    angular.element(document).ready(function () {
        $('#asnsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableASNsController.$inject = ['asns', '$scope', 'locationUtils'];
module.exports = TableASNsController;