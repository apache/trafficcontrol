/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 */

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
