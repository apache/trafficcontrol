var FormNewStatusController = function(status, $scope, $controller, statusService) {

    // extends the FormStatusController to inherit common methods
    angular.extend(this, $controller('FormStatusController', { status: status, $scope: $scope }));

    $scope.statusName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(status) {
        statusService.createStatus(status);
    };

};

FormNewStatusController.$inject = ['status', '$scope', '$controller', 'statusService'];
module.exports = FormNewStatusController;