var FormDeliveryServiceController = function(deliveryService, $scope, $uibModal, $anchorScroll, formUtils, stringUtils, locationUtils, cdnService, deliveryServiceService, profileService, typeService) {

    var deleteDeliveryService = function(ds) {
        deliveryServiceService.deleteDeliveryService(ds.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/delivery-services');
            });
    };

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

    $scope.deliveryServiceCopy = angular.copy(deliveryService);

    $scope.deliveryService = deliveryService;

    $scope.props = [
        { name: 'id', type: 'number', required: true, readonly: true },
        { name: 'displayName', type: 'text', required: true, maxLength: 48 },
        { name: 'xmlId', type: 'text', required: true, maxLength: 48 },
        { name: 'active', type: 'number', required: true, maxLength: 1 },
        { name: 'signed', type: 'number', required: true, maxLength: 1 },
        { name: 'qstringIgnore', type: 'number', required: true, maxLength: 1 },
        { name: 'geoLimit', type: 'number', required: true, maxLength: 1 },
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
        { name: 'ipv6RoutingEnabled', type: 'number', required: false, maxLength: 4 },
        { name: 'rangeRequestHandling', type: 'number', required: false, maxLength: 4 },
        { name: 'edgeHeaderRewrite', type: 'text', required: false, maxLength: 2048 },
        { name: 'midHeaderRewrite', type: 'text', required: false, maxLength: 2048 },
        { name: 'originShield', type: 'text', required: false, maxLength: 1024 },
        { name: 'regexRemap', type: 'text', required: false, maxLength: 1024 },
        { name: 'remapText', type: 'text', required: false, maxLength: 2048 },
        { name: 'cacheurl', type: 'text', required: false, maxLength: 1024 },
        { name: 'multiSiteOrigin', type: 'number', required: false, maxLength: 1 },
        { name: 'trResponseHeaders', type: 'text', required: false, maxLength: 1024 },
        { name: 'initialDispersion', type: 'number', required: false, maxLength: 11 },
        { name: 'dnsBypassCname', type: 'text', required: false, maxLength: 255 },
        { name: 'trRequestHeaders', type: 'text', required: false, maxLength: 1024 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.update = function(deliveryService) {
        deliveryServiceService.updateDeliveryService(deliveryService).
            then(function() {
                $scope.deliveryServiceCopy = angular.copy(deliveryService);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(ds) {
        var params = {
            title: 'Confirm Delete',
            message: 'This action CANNOT be undone. This will permanently delete ' + ds.displayName + '. Are you sure you want to delete ' + ds.displayName + '?'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deleteDeliveryService(ds);
        }, function () {
            // do nothing
        });
    };

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

FormDeliveryServiceController.$inject = ['deliveryService', '$scope', '$uibModal', '$anchorScroll', 'formUtils', 'stringUtils', 'locationUtils', 'cdnService', 'deliveryServiceService', 'profileService', 'typeService'];
module.exports = FormDeliveryServiceController;