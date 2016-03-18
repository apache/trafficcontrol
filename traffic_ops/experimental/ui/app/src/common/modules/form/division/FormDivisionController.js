var FormDivisionController = function(division, $scope, formUtils, stringUtils, locationUtils) {

    $scope.division = division;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormDivisionController.$inject = ['division', '$scope', 'formUtils', 'stringUtils', 'locationUtils'];
module.exports = FormDivisionController;