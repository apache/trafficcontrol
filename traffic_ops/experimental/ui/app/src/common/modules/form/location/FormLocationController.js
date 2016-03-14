var FormLocationController = function(location, $scope, formUtils, locationService) {

    $scope.location = location;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormLocationController.$inject = ['location', '$scope', 'formUtils', 'locationService'];
module.exports = FormLocationController;