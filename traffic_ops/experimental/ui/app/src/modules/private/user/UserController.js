var UserController = function($scope, $state, $uibModal, formUtils, locationUtils, userService, authService, roleService, userModel) {

    var updateUser = function(user, options) {
        userService.updateCurrentUser(user)
            .then(function() {
                if (options.signout) {
                    authService.logout();
                }
            });
    };

    var getRoles = function() {
        roleService.getRoles()
            .then(function(result) {
                $scope.roles = result;
            });
    };

    $scope.userName = angular.copy(userModel.user.username);

    $scope.user = userModel.user;

    $scope.confirmSave = function(user, usernameField) {
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

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getRoles();
    };
    init();

};

UserController.$inject = ['$scope', '$state', '$uibModal', 'formUtils', 'locationUtils', 'userService', 'authService', 'roleService', 'userModel'];
module.exports = UserController;
