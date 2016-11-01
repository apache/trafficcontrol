var FormDeliveryServiceController = function(deliveryService, $scope, formUtils, locationUtils, cdnService, profileService, typeService) {

    var getTypes = function() {
        typeService.getTypes('deliveryservice')
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

    $scope.falseTrue = [
        { value: false, label: 'false' },
        { value: true, label: 'true' }
    ];

    $scope.protocols = [
        { value: 0, label: '0 - HTTP' },
        { value: 1, label: '1 - HTTPS' },
        { value: 2, label: '2 - HTTP AND HTTPS' },
        { value: 3, label: '3 - HTTP TO HTTPS' }
    ];

    $scope.qStrings = [
        { value: 0, label: '0 - use qstring in cache key, and pass up' },
        { value: 1, label: '1 - ignore in cache key, and pass up' },
        { value: 2, label: '2 - drop at edge' }
    ];

    $scope.geoLimits = [
        { value: 0, label: '0 - None' },
        { value: 1, label: '1 - CZF only' },
        { value: 2, label: '2 - CZF + Country Code(s)' }
    ];

    $scope.geoProviders = [
        { value: 0, label: '0 - Maxmind (Default)' },
        { value: 1, label: '1 - Neustar' }
    ];

    $scope.dscps = [
        { value: 0, label: '0  - Best Effort' },
        { value: 10, label: '10 - AF11' },
        { value: 12, label: '12 - AF12' },
        { value: 14, label: '14 - AF13' },
        { value: 18, label: '18  - AF21' },
        { value: 20, label: '20  - AF22' },
        { value: 22, label: '22  - AF23' },
        { value: 26, label: '26  - AF31' },
        { value: 28, label: '28  - AF32' },
        { value: 30, label: '30  - AF33' },
        { value: 34, label: '34  - AF41' },
        { value: 36, label: '36  - AF42' },
        { value: 37, label: '37  - ' },
        { value: 38, label: '38  - AF43' },
        { value: 8, label: '8  - CS1' },
        { value: 16, label: '16  - CS2' },
        { value: 24, label: '24  - CS3' },
        { value: 32, label: '32  - CS4' },
        { value: 40, label: '40  - CS5' },
        { value: 48, label: '48  - CS6' },
        { value: 56, label: '56  - CS7' }
    ];

    $scope.dispersions = [
        { value: 1, label: '1 - OFF' },
        { value: 2, label: '2' },
        { value: 3, label: '3' },
        { value: 4, label: '4' },
        { value: 5, label: '5' },
        { value: 6, label: '6' },
        { value: 7, label: '7' },
        { value: 8, label: '8' },
        { value: 9, label: '9' },
        { value: 10, label: '10' }
    ];

    $scope.rrhs = [
        { value: 0, label: "0 - Don't cache" },
        { value: 1, label: "1 - Use background_fetch plugin" },
        { value: 2, label: "2 - Use cache_range_requests plugin" }
    ];

    $scope.msoAlgos = [
        { value: 0, label: "0 - Consistent Hash" },
        { value: 1, label: "1 - Primary/Backup" },
        { value: 2, label: "2 - Strict Round Robin" },
        { value: 3, label: "3 - IP-based Round Robin" },
        { value: 4, label: "4 - Latch on Failover" }
    ];

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