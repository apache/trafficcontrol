var FormUserController = function(user, $scope, $timeout, userService) {

    $scope.updateUser = function(user) {
        userService.updateCurrentUser(user);
    };

    $scope.resetUser = function() {
        $timeout(function() {
            $scope.userData = angular.copy(user.response[0]);
        });
    };
    $scope.resetUser();

};

FormUserController.$inject = ['user', '$scope', '$timeout', 'userService'];
module.exports = FormUserController;