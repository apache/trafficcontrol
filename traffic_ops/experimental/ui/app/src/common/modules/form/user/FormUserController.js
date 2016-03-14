var FormUserController = function(user, $scope, formUtils, userService) {

    var updateUser = function(user) {
        userService.updateUser(user);
    };

    $scope.userData = user;

    $scope.confirmUpdate = function(user, usernameField) {
        updateUser(user);
    };

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormUserController.$inject = ['user', '$scope', 'formUtils', 'userService'];
module.exports = FormUserController;