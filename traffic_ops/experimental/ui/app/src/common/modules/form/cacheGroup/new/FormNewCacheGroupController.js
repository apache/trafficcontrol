var FormNewCacheGroupController = function(cacheGroup, $scope, $controller, cacheGroupService) {

    // extends the FormCacheGroupController to inherit common methods
    angular.extend(this, $controller('FormCacheGroupController', { cacheGroup: cacheGroup, $scope: $scope }));

    $scope.cacheGroupName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(cacheGroup) {
        cacheGroupService.createCacheGroup(cacheGroup);
    };

};

FormNewCacheGroupController.$inject = ['cacheGroup', '$scope', '$controller', 'cacheGroupService'];
module.exports = FormNewCacheGroupController;