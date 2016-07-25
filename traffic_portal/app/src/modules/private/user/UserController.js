var UserController = function($scope, $state, $uibModal, $timeout, formUtils, deliveryServicesModel, userService, authService, userModel) {

    var updateUser = function(user, options) {
        user.token = null; // this will null out any token the user may have had
        userService.updateCurrentUser(user)
            .then(function() {
                if (options.signout) {
                    authService.logout();
                }
            });
    };

    $scope.deliveryServices = deliveryServicesModel.deliveryServices;

    $scope.showDS = function(dsId) {
        $state.go('trafficPortal.private.deliveryService.view.overview.detail', { deliveryServiceId: dsId } );
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
            });
        } else {
            updateUser(user, { signout : false });
        }
    };

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

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

UserController.$inject = ['$scope', '$state', '$uibModal', '$timeout', 'formUtils', 'deliveryServicesModel', 'userService', 'authService', 'userModel'];
module.exports = UserController;
