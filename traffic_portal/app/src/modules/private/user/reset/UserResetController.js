var UserResetController = function($scope, $timeout, userModel) {

    $scope.title = 'Reset User Password';

    $scope.reset = true;

    $scope.resetUser = function() {
        $timeout(function() {
            $scope.userData = angular.copy(userModel.user);
        });
    };
    $scope.resetUser();

    $scope.$on('userModel::userUpdated', function() {
        $scope.resetUser();
    });

};

UserResetController.$inject = ['$scope', '$timeout', 'userModel'];
module.exports = UserResetController;
