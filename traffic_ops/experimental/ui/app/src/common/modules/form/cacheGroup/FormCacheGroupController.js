var FormCacheGroupController = function(cacheGroup, $scope, cacheGroupService) {

    $scope.cacheGroup = cacheGroup;

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormCacheGroupController.$inject = ['cacheGroup', '$scope', 'cacheGroupService'];
module.exports = FormCacheGroupController;