var FormRegionController = function(region, $scope, formUtils, stringUtils, locationUtils, divisionService) {

    var getDivisions = function() {
        divisionService.getDivisions()
            .then(function(result) {
                $scope.divisions = result;
            });
    };

    $scope.region = region;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getDivisions();
    };
    init();

};

FormRegionController.$inject = ['region', '$scope', 'formUtils', 'stringUtils', 'locationUtils', 'divisionService'];
module.exports = FormRegionController;