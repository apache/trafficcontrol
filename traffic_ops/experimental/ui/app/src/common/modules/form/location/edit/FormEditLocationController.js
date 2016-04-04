var FormEditLocationController = function(location, $scope, $controller, $uibModal, $anchorScroll, locationUtils, locationService) {

    // extends the FormLocationController to inherit common methods
    angular.extend(this, $controller('FormLocationController', { location: location, $scope: $scope }));

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
            title: 'Delete Location: ' + location.name,
            key: location.name
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
            deleteLocation(location);
        }, function () {
            // do nothing
        });
    };

};

FormEditLocationController.$inject = ['location', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'locationService'];
module.exports = FormEditLocationController;