var FormNewCacheGroupController = function(cacheGroup, $scope, $controller, locationUtils, cacheGroupService) {

    // extends the FormCacheGroupController to inherit common methods
    angular.extend(this, $controller('FormCacheGroupController', { cacheGroup: cacheGroup, $scope: $scope }));

    $scope.cacheGroupName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(cacheGroup) {
        cacheGroupService.createCacheGroup(cacheGroup).
            then(function() {
                locationUtils.navigateToPath('/configure/cache-groups');
            });
    };

};

FormNewCacheGroupController.$inject = ['cacheGroup', '$scope', '$controller', 'locationUtils', 'cacheGroupService'];
module.exports = FormNewCacheGroupController;