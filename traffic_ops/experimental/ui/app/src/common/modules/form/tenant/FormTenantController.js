var FormTenantController = function(tenant, $scope, $uibModal, formUtils, locationUtils, tenantService) {

    var deleteTenant = function(tenant) {
        tenantService.deleteTenant(tenant.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/tenants');
            });
    };

    $scope.tenantOriginal = tenant;

    $scope.tenant = angular.copy(tenant);

    $scope.props = [
        { name: 'id', type: 'number', required: true, readonly: true },
        { name: 'name', type: 'text', required: true, maxLength: 45 }
    ];

    $scope.update = function(tenant) {
        alert('implement update');
    };

    $scope.confirmDelete = function(tenant) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + tenant.name + '. Are you sure you want to delete ' + tenant.name + '?'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
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

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormTenantController.$inject = ['tenant', '$scope', '$uibModal', 'formUtils', 'locationUtils', 'tenantService'];
module.exports = FormTenantController;