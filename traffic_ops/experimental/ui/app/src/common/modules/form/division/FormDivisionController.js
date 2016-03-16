var FormDivisionController = function(division, $scope, $uibModal, formUtils, locationUtils, divisionService) {

    var deleteDivision = function(division) {
        divisionService.deleteDivision(division.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/divisions');
            });
    };

    $scope.divisionOriginal = division;

    $scope.division = angular.copy(division);

    $scope.props = [
        { name: 'id', required: true, readonly: true },
        { name: 'name', required: true, maxLength: 45 }
    ];

    $scope.update = function(division) {
        alert('implement update');
    };

    $scope.confirmDelete = function(division) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + division.name + '. Are you sure you want to delete ' + division.name + '?'
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
            deleteDivision(division);
        }, function () {
            // do nothing
        });
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormDivisionController.$inject = ['division', '$scope', '$uibModal', 'formUtils', 'locationUtils', 'divisionService'];
module.exports = FormDivisionController;