var DialogConfirmController = function(params, $scope, $uibModalInstance) {

    $scope.params = params;

    $scope.yes = function() {
        $uibModalInstance.close();
    };

    $scope.no = function () {
        $uibModalInstance.dismiss('no');
    };

};

DialogConfirmController.$inject = ['params', '$scope', '$uibModalInstance'];
module.exports = DialogConfirmController;
