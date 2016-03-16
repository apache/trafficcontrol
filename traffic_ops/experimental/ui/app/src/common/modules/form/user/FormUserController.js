var FormUserController = function(user, showDelete, $scope, $uibModal, formUtils, locationUtils, roleService, userService) {

    var updateUser = function(user) {
        userService.updateUser(user);
    };

    var deleteUser = function(user) {
        userService.deleteUser(user.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/users');
            });
    };

    var getRoles = function() {
        roleService.getRoles()
            .then(function(result) {
                $scope.roles = result;
            });
    };

    $scope.userOriginal = angular.copy(user);

    $scope.user = user;

    $scope.showDelete = showDelete;

    $scope.confirmUpdate = function(user, usernameField) {
        updateUser(user);
    };

    $scope.confirmDelete = function(user) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + user.username + '. Are you sure you want to delete ' + user.username + '?'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deleteUser(user);
        }, function () {
            // do nothing
        });
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getRoles();
    };
    init();

};

FormUserController.$inject = ['user', 'showDelete', '$scope', '$uibModal', 'formUtils', 'locationUtils', 'roleService', 'userService'];
module.exports = FormUserController;