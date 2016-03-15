var FormDivisionController = function(division, $scope, formUtils, locationUtils, divisionService) {

    $scope.divisionOriginal = division;

    $scope.division = angular.copy(division);

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'name', required: true, maxLength: 45 }
    ];

    $scope.update = function(division) {
        alert('implement update');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormDivisionController.$inject = ['division', '$scope', 'formUtils', 'locationUtils', 'divisionService'];
module.exports = FormDivisionController;