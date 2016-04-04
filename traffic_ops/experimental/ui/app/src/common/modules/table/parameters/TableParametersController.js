var TableParametersController = function(parameters, $scope, locationUtils) {

    $scope.parameters = parameters;

    $scope.editParameter = function(id) {
        locationUtils.navigateToPath('/admin/parameters/' + id + '/edit');
    };

    $scope.createParameter = function() {
        locationUtils.navigateToPath('/admin/parameters/new');
    };

    angular.element(document).ready(function () {
        $('#parametersTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 100
        });
    });

};

TableParametersController.$inject = ['parameters', '$scope', 'locationUtils'];
module.exports = TableParametersController;