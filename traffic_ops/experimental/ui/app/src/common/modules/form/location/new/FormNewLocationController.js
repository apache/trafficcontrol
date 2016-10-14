var FormNewLocationController = function(location, $scope, $controller, locationService) {

    // extends the FormLocationController to inherit common methods
    angular.extend(this, $controller('FormLocationController', { location: location, $scope: $scope }));

    $scope.locationName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(location) {
        locationService.createLocation(location);
    };

};

FormNewLocationController.$inject = ['location', '$scope', '$controller', 'locationService'];
module.exports = FormNewLocationController;