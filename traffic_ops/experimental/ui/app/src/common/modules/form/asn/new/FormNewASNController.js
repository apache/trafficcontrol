var FormNewASNController = function(asn, $scope, $controller, locationUtils, asnService) {

    // extends the FormASNController to inherit common methods
    angular.extend(this, $controller('FormASNController', { asn: asn, $scope: $scope }));

    $scope.asnName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(asn) {
        asnService.createASN(asn).
            then(function() {
                locationUtils.navigateToPath('/admin/asns');
            });
    };

};

FormNewASNController.$inject = ['asn', '$scope', '$controller', 'locationUtils', 'asnService'];
module.exports = FormNewASNController;