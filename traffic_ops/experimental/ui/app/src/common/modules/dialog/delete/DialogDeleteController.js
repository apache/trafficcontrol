var DialogDeleteController = function(params, $scope, $uibModalInstance) {

    $scope.params = params;

    $scope.delete = function() {
        $uibModalInstance.close();
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };

};

DialogDeleteController.$inject = ['params', '$scope', '$uibModalInstance'];
module.exports = DialogDeleteController;
