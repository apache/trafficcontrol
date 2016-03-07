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

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormUserController.$inject = ['user', '$scope', '$timeout', 'userService'];
module.exports = FormUserController;