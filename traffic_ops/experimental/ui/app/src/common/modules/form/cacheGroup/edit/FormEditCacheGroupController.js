var FormEditCacheGroupController = function(cacheGroup, $scope, $controller, $uibModal, $anchorScroll, locationUtils, cacheGroupService) {

    // extends the FormCacheGroupController to inherit common methods
    angular.extend(this, $controller('FormCacheGroupController', { cacheGroup: cacheGroup, $scope: $scope }));

    var deleteCacheGroup = function(cacheGroup) {
        cacheGroupService.deleteCacheGroup(cacheGroup.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/cache-groups');
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
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + cacheGroup.name + '. Are you sure you want to delete ' + cacheGroup.name + '?'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
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