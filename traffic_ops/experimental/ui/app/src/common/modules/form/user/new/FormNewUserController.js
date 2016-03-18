var FormNewUserController = function(user, $scope, $controller, locationUtils, userService) {

    // extends the FormUserController to inherit common methods
    angular.extend(this, $controller('FormUserController', { user: user, $scope: $scope }));

    $scope.userName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.confirmSave = function(user, usernameField) {
        userService.createUser(user).
            then(function() {
                locationUtils.navigateToPath('/admin/users');
            });
    };

};

FormNewUserController.$inject = ['user', '$scope', '$controller', 'locationUtils', 'userService'];
module.exports = FormNewUserController;