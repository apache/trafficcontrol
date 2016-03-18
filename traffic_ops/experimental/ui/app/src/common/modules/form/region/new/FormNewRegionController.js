var FormNewRegionController = function(region, $scope, $controller, locationUtils, regionService) {

    // extends the FormRegionController to inherit common methods
    angular.extend(this, $controller('FormRegionController', { region: region, $scope: $scope }));

    $scope.regionName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(region) {
        regionService.createRegion(region).
            then(function() {
                locationUtils.navigateToPath('/admin/regions');
            });
    };

};

FormNewRegionController.$inject = ['region', '$scope', '$controller', 'locationUtils', 'regionService'];
module.exports = FormNewRegionController;