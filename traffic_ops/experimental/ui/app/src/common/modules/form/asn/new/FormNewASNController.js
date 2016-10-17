var FormNewASNController = function(asn, $scope, $controller, asnService) {

    // extends the FormASNController to inherit common methods
    angular.extend(this, $controller('FormASNController', { asn: asn, $scope: $scope }));

    $scope.asnName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(asn) {
        asnService.createASN(asn);
    };

};

FormNewASNController.$inject = ['asn', '$scope', '$controller', 'asnService'];
module.exports = FormNewASNController;