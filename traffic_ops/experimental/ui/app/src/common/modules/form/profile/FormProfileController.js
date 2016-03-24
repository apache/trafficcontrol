var FormProfileController = function(profile, $scope, formUtils, stringUtils, locationUtils) {

    $scope.profile = profile;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormProfileController.$inject = ['profile', '$scope', 'formUtils', 'stringUtils', 'locationUtils'];
module.exports = FormProfileController;