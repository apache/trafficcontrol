var FormRegionController = function(region, $scope, formUtils, regionService) {

    $scope.region = region;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormRegionController.$inject = ['region', '$scope', 'formUtils', 'regionService'];
module.exports = FormRegionController;