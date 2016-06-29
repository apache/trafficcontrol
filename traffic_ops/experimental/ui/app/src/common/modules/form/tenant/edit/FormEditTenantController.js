var FormEditTenantController = function(tenant, $scope, $controller, $uibModal, $anchorScroll, locationUtils, tenantService) {

    // extends the FormTenantController to inherit common methods
    angular.extend(this, $controller('FormTenantController', { tenant: tenant, $scope: $scope }));

    var deleteTenant = function(tenant) {
        tenantService.deleteTenant(tenant.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/tenants');
            });
    };

    $scope.tenantName = angular.copy(tenant.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(tenant) {
        tenantService.updateTenant(tenant).
            then(function() {
                $scope.tenantName = angular.copy(tenant.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(tenant) {
        var params = {
            title: 'Delete Tenant: ' + tenant.name,
            key: tenant.name
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
            controller: 'DialogDeleteController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deleteTenant(tenant);
        }, function () {
            // do nothing
        });
    };

};

FormEditTenantController.$inject = ['tenant', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'tenantService'];
module.exports = FormEditTenantController;