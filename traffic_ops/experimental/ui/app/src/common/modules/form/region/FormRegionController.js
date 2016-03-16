var FormRegionController = function(region, $scope, $uibModal, $anchorScroll, formUtils, locationUtils, divisionService, regionService) {

    var deleteRegion = function(region) {
        regionService.deleteRegion(region.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/regions');
            });
    };

    var getDivisions = function() {
        divisionService.getDivisions()
            .then(function(result) {
                $scope.divisions = result;
            });
    };

    $scope.regionCopy = angular.copy(region);

    $scope.region = region;

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'name', required: true, maxLength: 45 }
    ];

    $scope.update = function(region) {
        regionService.updateRegion(region).
            then(function() {
                $scope.regionCopy = angular.copy(region);
                $anchorScroll(); // scrolls window to top
            });
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

    var init = function () {
        getDivisions();
    };
    init();

};

FormRegionController.$inject = ['region', '$scope', '$uibModal', '$anchorScroll', 'formUtils', 'locationUtils', 'divisionService', 'regionService'];
module.exports = FormRegionController;