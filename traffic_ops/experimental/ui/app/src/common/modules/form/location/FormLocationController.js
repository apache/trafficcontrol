var FormLocationController = function(location, $scope, $uibModal, $anchorScroll, formUtils, locationUtils, locationService, regionService) {

    var deleteLocation = function(location) {
        locationService.deleteLocation(location.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/locations');
            });
    };

    var getRegions = function() {
        regionService.getRegions()
            .then(function(result) {
                $scope.regions = result;
            });
    };

    $scope.locationCopy = angular.copy(location);

    $scope.location = location;

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'name', required: true, maxLength: 45 },
        { name: 'shortName', required: true, maxLength: 12 },
        { name: 'address', required: true, maxLength: 128 },
        { name: 'city', required: true, maxLength: 128 },
        { name: 'state', required: true, maxLength: 2 },
        { name: 'zip', required: true, maxLength: 5 },
        { name: 'poc', required: false, maxLength: 128 },
        { name: 'phone', required: false, maxLength: 45 },
        { name: 'email', required: false, maxLength: 128 }
    ];

    $scope.update = function(location) {
        locationService.updateLocation(location).
            then(function() {
                $scope.locationCopy = angular.copy(location);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(location) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + location.name + '. Are you sure you want to delete ' + location.name + '?'
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
            deleteLocation(location);
        }, function () {
            // do nothing
        });
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getRegions();
    };
    init();

};

FormLocationController.$inject = ['location', '$scope', '$uibModal', '$anchorScroll', 'formUtils', 'locationUtils', 'locationService', 'regionService'];
module.exports = FormLocationController;