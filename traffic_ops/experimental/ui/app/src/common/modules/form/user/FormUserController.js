var FormUserController = function(user, $scope, formUtils, stringUtils, locationUtils, roleService) {

    var getRoles = function() {
        roleService.getRoles()
            .then(function(result) {
                $scope.roles = result;
            });
    };

    $scope.user = user;

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getRoles();
    };
    init();

};

FormUserController.$inject = ['user', '$scope', 'formUtils', 'stringUtils', 'locationUtils', 'roleService'];
module.exports = FormUserController;