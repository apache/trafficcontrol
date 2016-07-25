var HomeController = function($scope, $uibModal, propertiesModel, authService, userService) {

    $scope.sections = propertiesModel.properties.home.sections;

    $scope.credentials = {
        username: '',
        password: ''
    };

    $scope.login = function(credentials) {
        authService.login(credentials.username, credentials.password);
    };

    $scope.resetPassword = function() {

        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/reset/dialog.reset.tpl.html',
            controller: 'DialogResetController'
        });

        modalInstance.result.then(function(email) {
            userService.resetPassword(email);
        }, function () {
        });
    };

};

HomeController.$inject = ['$scope', '$uibModal', 'propertiesModel', 'authService', 'userService'];
module.exports = HomeController;