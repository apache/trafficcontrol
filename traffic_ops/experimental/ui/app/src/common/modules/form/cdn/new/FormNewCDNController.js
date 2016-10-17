var FormNewCDNController = function(cdn, $scope, $controller, cdnService) {

    // extends the FormCDNController to inherit common methods
    angular.extend(this, $controller('FormCDNController', { cdn: cdn, $scope: $scope }));

    $scope.cdnName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(cdn) {
        cdnService.createCDN(cdn)
    };

};

FormNewCDNController.$inject = ['cdn', '$scope', '$controller', 'cdnService'];
module.exports = FormNewCDNController;