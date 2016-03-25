var FormCDNController = function(cdn, $scope, formUtils, stringUtils, locationUtils) {

    $scope.cdn = cdn;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormCDNController.$inject = ['cdn', '$scope', 'formUtils', 'stringUtils', 'locationUtils'];
module.exports = FormCDNController;