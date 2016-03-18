var FormNewServerController = function(server, $scope, $controller, locationUtils, serverService) {

    // extends the FormServerController to inherit common methods
    angular.extend(this, $controller('FormServerController', { server: server, $scope: $scope }));

    $scope.serverName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(server) {
        serverService.createServer(server).
            then(function() {
                locationUtils.navigateToPath('/configure/servers');
            });
    };

};

FormNewServerController.$inject = ['server', '$scope', '$controller', 'locationUtils', 'serverService'];
module.exports = FormNewServerController;