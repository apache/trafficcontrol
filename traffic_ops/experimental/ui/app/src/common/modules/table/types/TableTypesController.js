var TableTypesController = function(types, $scope, locationUtils) {

    $scope.types = types;

    $scope.editType = function(id) {
        locationUtils.navigateToPath('/admin/types/' + id + '/edit');
    };

    $scope.createType = function() {
        locationUtils.navigateToPath('/admin/types/new');
    };

    angular.element(document).ready(function () {
        $('#typesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableTypesController.$inject = ['types', '$scope', 'locationUtils'];
module.exports = TableTypesController;