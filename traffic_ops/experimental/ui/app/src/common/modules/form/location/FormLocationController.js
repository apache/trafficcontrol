var FormLocationController = function(location, $scope, formUtils, stringUtils, locationUtils, regionService) {

    var getRegions = function() {
        regionService.getRegions()
            .then(function(result) {
                $scope.regions = result;
            });
    };

    $scope.location = location;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 },
        { name: 'shortName', type: 'text', required: true, maxLength: 12 },
        { name: 'address', type: 'text', required: true, maxLength: 128 },
        { name: 'city', type: 'text', required: true, maxLength: 128 },
        { name: 'state', type: 'text', required: true, maxLength: 2 },
        { name: 'zip', type: 'text', required: true, maxLength: 5 },
        { name: 'poc', type: 'text', required: false, maxLength: 128 },
        { name: 'phone', type: 'text', required: false, maxLength: 45 },
        { name: 'email', type: 'text', required: false, maxLength: 128 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getRegions();
    };
    init();

};

FormLocationController.$inject = ['location', '$scope', 'formUtils', 'stringUtils', 'locationUtils', 'regionService'];
module.exports = FormLocationController;