var FormLocationController = function(location, $scope, formUtils, locationUtils, locationService) {

    $scope.locationOriginal = location;

    $scope.location = angular.copy(location);

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'name', required: true, maxLength: 45 },
        { name: 'shortName', required: true, maxLength: 12 },
        { name: 'address', required: true, maxLength: 128 },
        { name: 'city', required: true, maxLength: 128 },
        { name: 'state', required: true, maxLength: 2 },
        { name: 'zip', required: true, maxLength: 5 },
        { name: 'poc', required: false, maxLength: 128 },
        { name: 'phone', required: false, maxLength: 45 },
        { name: 'email', required: false, maxLength: 128 }
    ];

    $scope.update = function(location) {
        alert('implement update');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormLocationController.$inject = ['location', '$scope', 'formUtils', 'locationUtils', 'locationService'];
module.exports = FormLocationController;