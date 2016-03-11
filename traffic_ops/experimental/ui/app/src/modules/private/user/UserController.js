var UserController = function($scope, $state, $uibModal, userService, authService, userModel) {

    var updateUser = function(user, options) {
        userService.updateCurrentUser(user)
            .then(function() {
                if (options.signout) {
                    authService.logout();
                }
            });
    };

    $scope.userData = userModel.user;

    $scope.confirmUpdate = function(user, usernameField) {
        if (usernameField.$dirty) {
            var params = {
                title: 'Reauthentication Required',
                message: 'Changing your username to ' + user.username + ' will require you to reauthenticate. Is that OK?'
            };
            var modalInstance = $uibModal.open({
                templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
                controller: 'DialogConfirmController',
                size: 'sm',
                resolve: {
                    params: function () {
                        return params;
                    }
                }
            });
            modalInstance.result.then(function() {
                updateUser(user, { signout : true });
            }, function () {
                // do nothing
            });
        } else {
            updateUser(user, { signout : false });
        }
    };

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

UserController.$inject = ['$scope', '$state', '$uibModal', 'userService', 'authService', 'userModel'];
module.exports = UserController;
