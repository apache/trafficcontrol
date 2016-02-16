var UserController = function($scope, $state, $uibModal, $timeout, userService, authService, userModel) {

    var updateUser = function(user, options) {
        user.token = null; // this will null out any token the user may have had
        userService.updateCurrentUser(user)
            .then(function() {
                if (options.signout) {
                    authService.logout();
                }
            });
    };

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
                $log.debug('Update user cancelled...');
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

UserController.$inject = ['$scope', '$state', '$uibModal', '$timeout', 'userService', 'authService', 'userModel'];
module.exports = UserController;
