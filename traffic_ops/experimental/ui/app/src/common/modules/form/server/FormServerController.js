var FormServerController = function(server, $scope, formUtils, serverService) {

    $scope.server = server;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormServerController.$inject = ['server', '$scope', 'formUtils', 'serverService'];
module.exports = FormServerController;