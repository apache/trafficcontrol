var CacheGroupsController = function(entityId, service, showDownload, $scope, $rootScope, $interval, jsonUtils, messageModel) {

    var locationsInterval;

    var getLocations = function(showTimeoutError) {
        if (!$rootScope.online) return;

        var ignoreLoadingBar = true;
        service.getCacheGroupHealth(entityId, ignoreLoadingBar, showTimeoutError)
            .then(
            function(response) {
                $scope.locationHealth = response;
            },
            function(fault) {
                // do nothing
            }).finally(function() {
                $scope.loaded = true;
            });
    };

    $scope.loaded = false;

    $scope.showDownload = showDownload;

    $scope.locationHealth = {
        totalOnline: 0,
        totalOffline: 0,
        locations: []
    };

    // pagination
    $scope.currentLocationPage = 1;
    $scope.locationsPerPage = 10;

    $scope.onlinePercent = function(location) {
        return location.online / (location.online + location.offline);
    };

    $scope.downloadCaches = function() {
        service.getServers(entityId, false)
            .then(
                function(response) {
                    jsonUtils.convertToCSV(response, 'Caches', ['hostName', 'domainName', 'type', 'cachegroup', 'ipAddress', 'ip6Address']);
                },
                function(fault) {
                    messageModel.setMessages([ { level: 'error', text: 'Failed to download cache servers.' } ], false);
                }
            );
    };

    $scope.$on("$destroy", function() {
        if (angular.isDefined(locationsInterval)) {
            $interval.cancel(locationsInterval);
            locationsInterval = undefined;
        }
    });

    var init = function () {
        getLocations(true);
        locationsInterval = $interval(function() { getLocations(false) }, 60 * 1000);
    };
    init();
};

CacheGroupsController.$inject = ['entityId', 'service', 'showDownload', '$scope', '$rootScope', '$interval', 'jsonUtils', 'messageModel'];
module.exports = CacheGroupsController;
