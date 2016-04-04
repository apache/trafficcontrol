var FormDeliveryServiceController = function(deliveryService, $scope, formUtils, stringUtils, locationUtils, cdnService, profileService, typeService) {

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

    $scope.props = [
        { name: 'displayName', type: 'text', required: true, maxLength: 48 },
        { name: 'xmlId', type: 'text', required: true, maxLength: 48 },
        { name: 'dscp', type: 'number', required: true, maxLength: 11 },
        { name: 'qstringIgnore', type: 'number', required: true, maxLength: 1 },
        { name: 'httpBypassFqdn', type: 'text', required: false, maxLength: 255 },
        { name: 'dnsBypassIp', type: 'text', required: false, maxLength: 45 },
        { name: 'dnsBypassIp6', type: 'text', required: false, maxLength: 45 },
        { name: 'dnsBypassTtl', type: 'number', required: false, maxLength: 11 },
        { name: 'orgServerFqdn', type: 'text', required: false, maxLength: 255 },
        { name: 'ccrDnsTtl', type: 'number', required: false, maxLength: 11 },
        { name: 'globalMaxMbps', type: 'number', required: false, maxLength: 11 },
        { name: 'globalMaxTps', type: 'number', required: false, maxLength: 11 },
        { name: 'longDesc', type: 'text', required: false, maxLength: 1024 },
        { name: 'longDesc1', type: 'text', required: false, maxLength: 1024 },
        { name: 'longDesc2', type: 'text', required: false, maxLength: 1024 },
        { name: 'maxDnsAnswers', type: 'number', required: false, maxLength: 11 },
        { name: 'infoUrl', type: 'text', required: false, maxLength: 255 },
        { name: 'missLat', type: 'number', required: false, maxLength: 255, pattern: new RegExp('^[-+]?[0-9]*\.?[0-9]+$'), invalidMsg: 'Invalid coordinate' },
        { name: 'missLong', type: 'number', required: false, maxLength: 255, pattern: new RegExp('^[-+]?[0-9]*\.?[0-9]+$'), invalidMsg: 'Invalid coordinate' },
        { name: 'checkPath', type: 'text', required: false, maxLength: 255 },
        { name: 'protocol', type: 'number', required: false, maxLength: 4 },
        { name: 'sslKeyVersion', type: 'number', required: false, maxLength: 11 },
        { name: 'rangeRequestHandling', type: 'number', required: false, maxLength: 4 },
        { name: 'edgeHeaderRewrite', type: 'text', required: false, maxLength: 2048 },
        { name: 'midHeaderRewrite', type: 'text', required: false, maxLength: 2048 },
        { name: 'originShield', type: 'text', required: false, maxLength: 1024 },
        { name: 'regexRemap', type: 'text', required: false, maxLength: 1024 },
        { name: 'remapText', type: 'text', required: false, maxLength: 2048 },
        { name: 'cacheurl', type: 'text', required: false, maxLength: 1024 },
        { name: 'trResponseHeaders', type: 'text', required: false, maxLength: 1024 },
        { name: 'initialDispersion', type: 'number', required: false, maxLength: 11 },
        { name: 'dnsBypassCname', type: 'text', required: false, maxLength: 255 },
        { name: 'trRequestHeaders', type: 'text', required: false, maxLength: 1024 }
    ];

    $scope.labelize = stringUtils.labelize;

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

FormDeliveryServiceController.$inject = ['deliveryService', '$scope', 'formUtils', 'stringUtils', 'locationUtils', 'cdnService', 'profileService', 'typeService'];
module.exports = FormDeliveryServiceController;