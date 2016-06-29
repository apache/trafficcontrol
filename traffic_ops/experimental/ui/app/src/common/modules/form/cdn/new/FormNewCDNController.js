var FormNewCDNController = function(cdn, $scope, $controller, locationUtils, cdnService) {

    // extends the FormCDNController to inherit common methods
    angular.extend(this, $controller('FormCDNController', { cdn: cdn, $scope: $scope }));

    $scope.cdnName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(cdn) {
        cdnService.createCDN(cdn).
            then(function() {
                locationUtils.navigateToPath('/admin/cdns');
            });
    };

};

FormNewCDNController.$inject = ['cdn', '$scope', '$controller', 'locationUtils', 'cdnService'];
module.exports = FormNewCDNController;