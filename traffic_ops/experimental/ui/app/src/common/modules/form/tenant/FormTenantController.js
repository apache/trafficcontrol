var FormTenantController = function(tenant, $scope, formUtils, stringUtils, locationUtils) {

    $scope.tenant = angular.copy(tenant);

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormTenantController.$inject = ['tenant', '$scope', 'formUtils', 'stringUtils', 'locationUtils'];
module.exports = FormTenantController;