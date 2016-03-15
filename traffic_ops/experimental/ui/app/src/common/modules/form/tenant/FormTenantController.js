var FormTenantController = function(tenant, $scope, formUtils, locationUtils, tenantService) {

    $scope.tenantOriginal = tenant;

    $scope.tenant = angular.copy(tenant);

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'name', required: true, maxLength: 45 }
    ];

    $scope.update = function(tenant) {
        alert('implement update');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormTenantController.$inject = ['tenant', '$scope', 'formUtils', 'locationUtils', 'tenantService'];
module.exports = FormTenantController;