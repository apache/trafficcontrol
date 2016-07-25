var ReleaseVersionController = function(params, $scope, $uibModalInstance) {

    $scope.params = params;

    $scope.dismiss = function () {
        $uibModalInstance.dismiss('cancel');
    };

};

ReleaseVersionController.$inject = ['params', '$scope', '$uibModalInstance'];
module.exports = ReleaseVersionController;
