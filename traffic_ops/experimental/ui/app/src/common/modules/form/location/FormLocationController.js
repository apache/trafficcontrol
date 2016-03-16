var FormLocationController = function(location, $scope, $uibModal, formUtils, locationUtils, locationService) {

    var deleteLocation = function(location) {
        locationService.deleteLocation(location.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/locations');
            });
    };

    $scope.locationOriginal = location;

    $scope.location = angular.copy(location);

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
        alert('implement update');
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

};

FormLocationController.$inject = ['location', '$scope', '$uibModal', 'formUtils', 'locationUtils', 'locationService'];
module.exports = FormLocationController;