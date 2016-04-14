var FormNewParameterController = function(parameter, $scope, $controller, locationUtils, parameterService) {

    // extends the FormParameterController to inherit common methods
    angular.extend(this, $controller('FormParameterController', { parameter: parameter, $scope: $scope }));

    $scope.parameterName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(parameter) {
        parameterService.createParameter(parameter).
            then(function() {
                locationUtils.navigateToPath('/admin/parameters');
            });
    };

};

FormNewParameterController.$inject = ['parameter', '$scope', '$controller', 'locationUtils', 'parameterService'];
module.exports = FormNewParameterController;