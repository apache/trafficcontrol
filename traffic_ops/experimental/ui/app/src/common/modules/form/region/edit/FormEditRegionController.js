var FormEditRegionController = function(region, $scope, $controller, $uibModal, $anchorScroll, locationUtils, regionService) {

    // extends the FormRegionController to inherit common methods
    angular.extend(this, $controller('FormRegionController', { region: region, $scope: $scope }));

    var deleteRegion = function(region) {
        regionService.deleteRegion(region.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/regions');
            });
    };

    $scope.regionName = angular.copy(region.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(region) {
        regionService.updateRegion(region).
            then(function() {
                $scope.regionName = angular.copy(region.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(region) {
        var params = {
            title: 'Delete Region: ' + region.name,
            key: region.name
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
            deleteRegion(region);
        }, function () {
            // do nothing
        });
    };

};

FormEditRegionController.$inject = ['region', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'regionService'];
module.exports = FormEditRegionController;