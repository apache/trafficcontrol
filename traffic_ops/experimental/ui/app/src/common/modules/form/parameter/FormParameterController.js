var FormParameterController = function(parameter, $scope, formUtils, stringUtils, locationUtils) {

    $scope.parameter = parameter;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 1024 },
        { name: 'configFile', type: 'text', required: true, maxLength: 45 },
        { name: 'value', type: 'text', required: true, maxLength: 1024 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormParameterController.$inject = ['parameter', '$scope', 'formUtils', 'stringUtils', 'locationUtils'];
module.exports = FormParameterController;