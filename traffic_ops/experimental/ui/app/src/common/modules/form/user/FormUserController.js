var FormUserController = function(user, $scope, userService) {

    var updateUser = function(user) {
        userService.updateUser(user);
    };

    $scope.userData = user;

    $scope.confirmUpdate = function(user, usernameField) {
        updateUser(user);
    };

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormUserController.$inject = ['user', '$scope', 'userService'];
module.exports = FormUserController;