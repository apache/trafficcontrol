var FormEditDivisionController = function(division, $scope, $controller, $uibModal, $anchorScroll, locationUtils, divisionService) {

    // extends the FormDivisionController to inherit common methods
    angular.extend(this, $controller('FormDivisionController', { division: division, $scope: $scope }));

    var deleteDivision = function(division) {
        divisionService.deleteDivision(division.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/divisions');
            });
    };

    $scope.divisionName = angular.copy(division.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(division) {
        divisionService.updateDivision(division).
            then(function() {
                $scope.divisionName = angular.copy(division.name);
                $anchorScroll(); // scrolls window to top
            });
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

};

FormEditDivisionController.$inject = ['division', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'divisionService'];
module.exports = FormEditDivisionController;