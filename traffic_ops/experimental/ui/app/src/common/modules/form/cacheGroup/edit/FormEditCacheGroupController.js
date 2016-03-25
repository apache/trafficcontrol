var FormEditCacheGroupController = function(cacheGroup, $scope, $controller, $uibModal, $anchorScroll, locationUtils, cacheGroupService) {

    // extends the FormCacheGroupController to inherit common methods
    angular.extend(this, $controller('FormCacheGroupController', { cacheGroup: cacheGroup, $scope: $scope }));

    var deleteCacheGroup = function(cacheGroup) {
        cacheGroupService.deleteCacheGroup(cacheGroup.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/cache-groups');
            });
    };

    $scope.cacheGroupName = angular.copy(cacheGroup.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(cacheGroup) {
        cacheGroupService.updateCacheGroup(cacheGroup).
            then(function() {
                $scope.cacheGroupName = angular.copy(cacheGroup.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(cacheGroup) {
        var params = {
            title: 'Delete Cache Group: ' + cacheGroup.name,
            key: cacheGroup.name
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
            controller: 'DialogDeleteController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deleteCacheGroup(cacheGroup);
        }, function () {
            // do nothing
        });
    };

};

FormEditCacheGroupController.$inject = ['cacheGroup', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'cacheGroupService'];
module.exports = FormEditCacheGroupController;