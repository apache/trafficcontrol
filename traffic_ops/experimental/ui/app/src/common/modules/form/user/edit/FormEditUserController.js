var FormEditUserController = function(user, $scope, $controller, $uibModal, $anchorScroll, locationUtils, userService) {

    // extends the FormUserController to inherit common methods
    angular.extend(this, $controller('FormUserController', { user: user, $scope: $scope }));

    var saveUser = function(user) {
        userService.updateUser(user).
            then(function() {
                $scope.userName = angular.copy(user.username);
                $anchorScroll(); // scrolls window to top
            });
    };

    var deleteUser = function(user) {
        userService.deleteUser(user.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/users');
            });
    };

    $scope.userName = angular.copy(user.username);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.confirmSave = function(user, usernameField) {
        saveUser(user);
    };

    $scope.confirmDelete = function(user) {
        var params = {
            title: 'Delete User: ' + user.username,
            key: user.username
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
            controller: 'DialogDeleteController',
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

};

FormEditUserController.$inject = ['user', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'userService'];
module.exports = FormEditUserController;