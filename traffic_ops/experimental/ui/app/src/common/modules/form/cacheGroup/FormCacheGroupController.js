var FormCacheGroupController = function(cacheGroup, $scope, formUtils, locationUtils, cacheGroupService, typeService) {

    var getCacheGroups = function() {
        cacheGroupService.getCacheGroups()
            .then(function(result) {
                $scope.cacheGroups = result;
            });
    };

    var getTypes = function() {
        typeService.getTypes('cachegroup')
            .then(function(result) {
                $scope.types = result;
            });
    };

    $scope.cacheGroup = cacheGroup;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getCacheGroups();
        getTypes();
    };
    init();

};

FormCacheGroupController.$inject = ['cacheGroup', '$scope', 'formUtils', 'locationUtils', 'cacheGroupService', 'typeService'];
module.exports = FormCacheGroupController;