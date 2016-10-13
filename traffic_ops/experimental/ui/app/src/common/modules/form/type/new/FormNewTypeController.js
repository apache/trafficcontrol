var FormNewTypeController = function(type, $scope, $controller, typeService) {

    // extends the FormTypeController to inherit common methods
    angular.extend(this, $controller('FormTypeController', { type: type, $scope: $scope }));

    $scope.typeName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(type) {
        typeService.createType(type);
    };

};

FormNewTypeController.$inject = ['type', '$scope', '$controller', 'typeService'];
module.exports = FormNewTypeController;