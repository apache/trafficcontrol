var FormNewDivisionController = function(division, $scope, $controller, divisionService) {

    // extends the FormDivisionController to inherit common methods
    angular.extend(this, $controller('FormDivisionController', { division: division, $scope: $scope }));

    $scope.divisionName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(division) {
        divisionService.createDivision(division)
    };

};

FormNewDivisionController.$inject = ['division', '$scope', '$controller', 'divisionService'];
module.exports = FormNewDivisionController;