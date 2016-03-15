var FormRegionController = function(region, $scope, formUtils, locationUtils, regionService) {

    $scope.regionOriginal = region;

    $scope.region = angular.copy(region);

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'name', required: true, maxLength: 45 }
    ];

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormRegionController.$inject = ['region', '$scope', 'formUtils', 'locationUtils', 'regionService'];
module.exports = FormRegionController;