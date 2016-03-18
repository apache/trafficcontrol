var FormCacheGroupController = function(cacheGroup, $scope, formUtils, stringUtils, locationUtils, cacheGroupService, typeService) {

    var getCacheGroups = function() {
        cacheGroupService.getCacheGroups()
            .then(function(result) {
                $scope.cacheGroups = result;
            });
    };

    var getTypes = function() {
        typeService.getTypes()
            .then(function(result) {
                $scope.types = result;
            });
    };

    $scope.cacheGroup = cacheGroup;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 },
        { name: 'shortName', type: 'text', required: true, maxLength: 255 },
        { name: 'latitude', type: 'number', required: false, pattern: new RegExp('^[-+]?[0-9]*\.?[0-9]+$'), invalidMsg: 'Invalid coordinate' },
        { name: 'longitude', type: 'number', required: false, pattern: new RegExp('^[-+]?[0-9]*\.?[0-9]+$'), invalidMsg: 'Invalid coordinate' }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getCacheGroups();
        getTypes();
    };
    init();

};

FormCacheGroupController.$inject = ['cacheGroup', '$scope', 'formUtils', 'stringUtils', 'locationUtils', 'cacheGroupService', 'typeService'];
module.exports = FormCacheGroupController;