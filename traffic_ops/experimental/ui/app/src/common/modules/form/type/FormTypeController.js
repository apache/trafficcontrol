var FormTypeController = function(type, $scope, formUtils, stringUtils, locationUtils) {

    $scope.type = type;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 },
        { name: 'useInTable', type: 'text', required: true, maxLength: 45 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormTypeController.$inject = ['type', '$scope', 'formUtils', 'stringUtils', 'locationUtils'];
module.exports = FormTypeController;