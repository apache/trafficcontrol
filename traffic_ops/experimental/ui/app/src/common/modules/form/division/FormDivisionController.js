var FormDivisionController = function(division, $scope, divisionService) {

    $scope.division = division;

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormDivisionController.$inject = ['division', '$scope', 'divisionService'];
module.exports = FormDivisionController;