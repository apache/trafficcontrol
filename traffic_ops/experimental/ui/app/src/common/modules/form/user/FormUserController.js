var FormUserController = function(user, $scope, $timeout, userService, ENV) {

    var updateUser = function(user) {
        userService.updateUser(ENV.api['root'] + 'tm_user/' + user.id, user);
    };

    $scope.confirmUpdate = function(user, usernameField) {
        updateUser(user);
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

FormUserController.$inject = ['user', '$scope', '$timeout', 'userService', 'ENV'];
module.exports = FormUserController;