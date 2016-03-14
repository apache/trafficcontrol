var FormCacheGroupController = function(cacheGroup, $scope, formUtils, cacheGroupService) {

    $scope.cacheGroup = cacheGroup;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormCacheGroupController.$inject = ['cacheGroup', '$scope', 'formUtils', 'cacheGroupService'];
module.exports = FormCacheGroupController;