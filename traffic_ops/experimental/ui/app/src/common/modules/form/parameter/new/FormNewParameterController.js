var FormNewParameterController = function(parameter, $scope, $controller, parameterService) {

    // extends the FormParameterController to inherit common methods
    angular.extend(this, $controller('FormParameterController', { parameter: parameter, $scope: $scope }));

    $scope.parameterName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(parameter) {
        parameterService.createParameter(parameter);
    };

};

FormNewParameterController.$inject = ['parameter', '$scope', '$controller', 'parameterService'];
module.exports = FormNewParameterController;