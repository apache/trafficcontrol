var FormUserController = function(user, $scope, formUtils, locationUtils, userService) {

    var updateUser = function(user) {
        userService.updateUser(user);
    };

    $scope.userOriginal = angular.copy(user);

    $scope.user = user;

    $scope.confirmUpdate = function(user, usernameField) {
        updateUser(user);
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormUserController.$inject = ['user', '$scope', 'formUtils', 'locationUtils', 'userService'];
module.exports = FormUserController;