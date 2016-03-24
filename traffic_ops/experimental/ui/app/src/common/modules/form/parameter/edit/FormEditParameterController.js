var FormEditParameterController = function(parameter, $scope, $controller, $uibModal, $anchorScroll, locationUtils, parameterService) {

    // extends the FormParameterController to inherit common methods
    angular.extend(this, $controller('FormParameterController', { parameter: parameter, $scope: $scope }));

    var deleteParameter = function(parameter) {
        parameterService.deleteParameter(parameter.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/parameters');
            });
    };

    $scope.parameterName = angular.copy(parameter.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(parameter) {
        parameterService.updateParameter(parameter).
            then(function() {
                $scope.parameterName = angular.copy(parameter.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(parameter) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + parameter.name + '. Are you sure you want to delete ' + parameter.name + '?'
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
            deleteParameter(parameter);
        }, function () {
            // do nothing
        });
    };

};

FormEditParameterController.$inject = ['parameter', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'parameterService'];
module.exports = FormEditParameterController;