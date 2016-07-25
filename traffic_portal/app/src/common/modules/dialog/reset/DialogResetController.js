var DialogResetController = function($scope, $uibModalInstance, formUtils) {

    $scope.userData = {
        email: ""
    };

    $scope.reset = function (email) {
        $uibModalInstance.close(email);
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

DialogResetController.$inject = ['$scope', '$uibModalInstance', 'formUtils'];
module.exports = DialogResetController;
