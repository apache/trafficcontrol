var DialogResetController = function($scope, $uibModalInstance) {

    $scope.userData = {
        email: ""
    };

    $scope.reset = function (email) {
        $uibModalInstance.close(email);
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

DialogResetController.$inject = ['$scope', '$uibModalInstance'];
module.exports = DialogResetController;
