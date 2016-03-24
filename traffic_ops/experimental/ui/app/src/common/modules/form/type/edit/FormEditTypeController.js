var FormEditTypeController = function(type, $scope, $controller, $uibModal, $anchorScroll, locationUtils, typeService) {

    // extends the FormTypeController to inherit common methods
    angular.extend(this, $controller('FormTypeController', { type: type, $scope: $scope }));

    var deleteType = function(type) {
        typeService.deleteType(type.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/types');
            });
    };

    $scope.typeName = angular.copy(type.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(type) {
        typeService.updateType(type).
            then(function() {
                $scope.typeName = angular.copy(type.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(type) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + type.name + '. Are you sure you want to delete ' + type.name + '?'
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
            deleteType(type);
        }, function () {
            // do nothing
        });
    };

};

FormEditTypeController.$inject = ['type', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'typeService'];
module.exports = FormEditTypeController;