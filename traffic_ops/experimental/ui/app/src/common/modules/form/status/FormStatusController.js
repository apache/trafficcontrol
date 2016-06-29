var FormStatusController = function(status, $scope, formUtils, stringUtils, locationUtils) {

    $scope.status = status;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormStatusController.$inject = ['status', '$scope', 'formUtils', 'stringUtils', 'locationUtils'];
module.exports = FormStatusController;