var FormCacheGroupController = function(cacheGroup, $scope, $uibModal, formUtils, locationUtils, cacheGroupService) {

    var deleteCacheGroup = function(cacheGroup) {
        cacheGroupService.deleteCacheGroup(cacheGroup.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/cache-groups');
            });
    };

    $scope.cacheGroupOriginal = cacheGroup;

    $scope.cacheGroup = angular.copy(cacheGroup);

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'name', required: true, maxLength: 45 }
    ];

    $scope.update = function(cacheGroup) {
        alert('implement update');
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

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormCacheGroupController.$inject = ['cacheGroup', '$scope', '$uibModal', 'formUtils', 'locationUtils', 'cacheGroupService'];
module.exports = FormCacheGroupController;