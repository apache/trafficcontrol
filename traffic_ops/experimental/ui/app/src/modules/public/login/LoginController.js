var LoginController = function($scope, $log, $uibModal, authService) {

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
            authService.resetPassword(email);
        }, function () {
        });
    };

    var init = function() {};
    init();
};

LoginController.$inject = ['$scope', '$log', '$uibModal', 'authService'];
module.exports = LoginController;
