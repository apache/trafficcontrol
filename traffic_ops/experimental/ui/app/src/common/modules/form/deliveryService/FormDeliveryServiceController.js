var FormDeliveryServiceController = function(deliveryService, $scope, formUtils, locationUtils, cdnService, profileService, typeService) {

    var getTypes = function() {
        typeService.getTypes()
            .then(function(result) {
                $scope.types = result;
            });
    };

    var getCDNs = function() {
        cdnService.getCDNs()
            .then(function(result) {
                $scope.cdns = result;
            });
    };

    var getProfiles = function() {
        profileService.getProfiles()
            .then(function(result) {
                $scope.profiles = result;
            });
    };

    $scope.deliveryService = deliveryService;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getTypes();
        getCDNs();
        getProfiles();
    };
    init();

};

FormDeliveryServiceController.$inject = ['deliveryService', '$scope', 'formUtils', 'locationUtils', 'cdnService', 'profileService', 'typeService'];
module.exports = FormDeliveryServiceController;