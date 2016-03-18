var FormNewDivisionController = function(division, $scope, $controller, locationUtils, divisionService) {

    // extends the FormDivisionController to inherit common methods
    angular.extend(this, $controller('FormDivisionController', { division: division, $scope: $scope }));

    $scope.divisionName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(division) {
        divisionService.createDivision(division).
            then(function() {
                locationUtils.navigateToPath('/admin/divisions');
            });
    };

};

FormNewDivisionController.$inject = ['division', '$scope', '$controller', 'locationUtils', 'divisionService'];
module.exports = FormNewDivisionController;