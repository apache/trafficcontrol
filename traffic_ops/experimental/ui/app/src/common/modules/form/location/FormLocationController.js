var FormLocationController = function(location, $scope, locationService) {

    $scope.location = location;

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormLocationController.$inject = ['location', '$scope', 'locationService'];
module.exports = FormLocationController;