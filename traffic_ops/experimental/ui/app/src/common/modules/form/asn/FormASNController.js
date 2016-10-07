var FormASNController = function(asn, $scope, formUtils, locationUtils, cacheGroupService) {

    var getCacheGroups = function() {
        cacheGroupService.getCacheGroups()
            .then(function(result) {
                $scope.cachegroups = result;
            });
    };

    $scope.asn = asn;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getCacheGroups();
    };
    init();

};

FormASNController.$inject = ['asn', '$scope', 'formUtils', 'locationUtils', 'cacheGroupService'];
module.exports = FormASNController;