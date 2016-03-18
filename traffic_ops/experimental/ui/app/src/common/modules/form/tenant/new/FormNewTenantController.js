var FormNewTenantController = function(tenant, $scope, $controller, locationUtils, tenantService) {

    // extends the FormTenantController to inherit common methods
    angular.extend(this, $controller('FormTenantController', { tenant: tenant, $scope: $scope }));

    $scope.tenantName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(tenant) {
        tenantService.createTenant(tenant).
            then(function() {
                locationUtils.navigateToPath('/admin/tenants');
            });
    };

};

FormNewTenantController.$inject = ['tenant', '$scope', '$controller', 'locationUtils', 'tenantService'];
module.exports = FormNewTenantController;