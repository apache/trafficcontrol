var FormEditLocationController = function(location, $scope, $controller, $uibModal, $anchorScroll, locationUtils, locationService) {

    // extends the FormLocationController to inherit common methods
    angular.extend(this, $controller('FormLocationController', { location: location, $scope: $scope }));
//    debugger;

    var deleteLocation = function(location) {
        locationService.deleteLocation(location.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/locations');
            });
    };

    $scope.locationName = angular.copy(location.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(location) {
        locationService.updateLocation(location).
            then(function() {
                $scope.locationName = angular.copy(location.name);
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

};

FormEditLocationController.$inject = ['location', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'locationService'];
module.exports = FormEditLocationController;