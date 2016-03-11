var FormServerController = function(server, $scope, serverService) {

    $scope.server = server;

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormServerController.$inject = ['server', '$scope', 'serverService'];
module.exports = FormServerController;