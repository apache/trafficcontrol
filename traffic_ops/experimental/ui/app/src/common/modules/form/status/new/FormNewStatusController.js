var FormNewStatusController = function(status, $scope, $controller, locationUtils, statusService) {

    // extends the FormStatusController to inherit common methods
    angular.extend(this, $controller('FormStatusController', { status: status, $scope: $scope }));

    $scope.statusName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(status) {
        statusService.createStatus(status).
            then(function() {
                locationUtils.navigateToPath('/admin/statuses');
            });
    };

};

FormNewStatusController.$inject = ['status', '$scope', '$controller', 'locationUtils', 'statusService'];
module.exports = FormNewStatusController;