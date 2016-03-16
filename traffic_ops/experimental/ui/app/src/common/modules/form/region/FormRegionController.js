var FormRegionController = function(region, $scope, $uibModal, formUtils, locationUtils, regionService) {

    var deleteRegion = function(region) {
        regionService.deleteRegion(region.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/regions');
            });
    };

    $scope.regionOriginal = region;

    $scope.region = angular.copy(region);

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'name', required: true, maxLength: 45 }
    ];

    $scope.update = function(region) {
        alert('implement update');
    };

    $scope.confirmDelete = function(region) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + region.name + '. Are you sure you want to delete ' + region.name + '?'
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
            deleteRegion(region);
        }, function () {
            // do nothing
        });
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormRegionController.$inject = ['region', '$scope', '$uibModal', 'formUtils', 'locationUtils', 'regionService'];
module.exports = FormRegionController;