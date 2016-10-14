var FormNewServerController = function(server, $scope, $controller, serverService) {

    // extends the FormServerController to inherit common methods
    angular.extend(this, $controller('FormServerController', { server: server, $scope: $scope }));

    $scope.serverName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(server) {
        serverService.createServer(server);
    };

};

FormNewServerController.$inject = ['server', '$scope', '$controller', 'serverService'];
module.exports = FormNewServerController;