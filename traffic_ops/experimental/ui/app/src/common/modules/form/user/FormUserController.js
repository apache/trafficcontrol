var FormUserController = function(user, $scope, userService) {

    $scope.userData = angular.copy(user);

    $scope.updateUser = function(user) {
        userService.updateCurrentUser(user);
    };

};

FormUserController.$inject = ['user', '$scope', 'userService'];
module.exports = FormUserController;