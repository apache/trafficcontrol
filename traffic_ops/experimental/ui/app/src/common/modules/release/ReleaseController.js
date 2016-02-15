var ReleaseController = function(params, $scope, $uibModalInstance) {

    $scope.params = params;

    $scope.dismiss = function () {
        $uibModalInstance.dismiss('cancel');
    };

};

ReleaseController.$inject = ['params', '$scope', '$uibModalInstance'];
module.exports = ReleaseController;
