var FormASNController = function(asn, $scope, formUtils, stringUtils, locationUtils, cacheGroupService) {

    var getCacheGroups = function() {
        cacheGroupService.getCacheGroups()
            .then(function(result) {
                $scope.cachegroups = result;
            });
    };

    $scope.asn = asn;

    $scope.props = [
        { name: 'asn', type: 'number', required: true, maxLength: 11 },
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getCacheGroups();
    };
    init();

};

FormASNController.$inject = ['asn', '$scope', 'formUtils', 'stringUtils', 'locationUtils', 'cacheGroupService'];
module.exports = FormASNController;