var FormNewTypeController = function(type, $scope, $controller, locationUtils, typeService) {

    // extends the FormTypeController to inherit common methods
    angular.extend(this, $controller('FormTypeController', { type: type, $scope: $scope }));

    $scope.typeName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(type) {
        typeService.createType(type).
            then(function() {
                locationUtils.navigateToPath('/admin/types');
            });
    };

};

FormNewTypeController.$inject = ['type', '$scope', '$controller', 'locationUtils', 'typeService'];
module.exports = FormNewTypeController;