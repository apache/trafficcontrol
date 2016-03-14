var FormDivisionController = function(division, $scope, locationUtils, divisionService) {

    $scope.divisionOriginal = division;

    $scope.division = angular.copy(division);

    $scope.update = function(division) {
        alert('implement update');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormDivisionController.$inject = ['division', '$scope', 'locationUtils', 'divisionService'];
module.exports = FormDivisionController;