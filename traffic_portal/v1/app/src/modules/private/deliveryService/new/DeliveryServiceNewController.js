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

var DeliveryServiceNewController = function($scope, $location, $anchorScroll, formUtils, deliveryServiceService) {

    $scope.templates = {
        signedUrls: 'signedUrl.html',
        queryStringHandling: 'queryStringHandling.html',
        rangeRequestHandling: 'rangeRequestHandling.html',
        headerRewriteEdge: 'headerRewriteEdge.html',
        headerRewriteMid: 'headerRewriteMid.html',
        headerRewriteRedirectRouter: 'headerRewriteRedirectRouter.html'
    };

    $scope.openTab = 'service-desc-tab';

    $scope.booleans = [
        { label: 'Yes', value: true },
        { label: 'No', value: false }
    ];

    $scope.serviceDescStarted = true;
    $scope.trafficStarted = false;
    $scope.originStarted = false;
    $scope.coreStarted = false;
    $scope.serviceLimitStarted = false;
    $scope.headerStarted = false;
    $scope.notesStarted = false;
    $scope.confirmStarted = false;

    $scope.resetDSData = function() {
        $scope.dsData = {
            customer: '',
            contentType: 'video-on-demand',
            deliveryProtocol: 'http',
            routingType: 'dns',
            serviceDesc: '',
            peakBPSEstimate: 'less-than-5-Gbps',
            peakTPSEstimate: 'less-than-1000-TPS',
            maxLibrarySizeEstimate: 'less-than-200-GB',
            originURL: '',
            hasOriginDynamicRemap: false,
            originTestFile: '',
            hasOriginACLWhitelist: false,
            originHeaders: '',
            otherOriginSecurity: '',
            queryStringHandling: 'ignore-in-cache-key-and-pass-up',
            rangeRequestHandling: 'range-requests-not-used',
            hasSignedURLs: false,
            hasNegativeCachingCustomization: false,
            negativeCachingCustomizationNote: '',
            serviceAliases: [ '' ],
            rateLimitingGBPS: 0,
            rateLimitingTPS: 0,
            overflowService: '',
            headerRewriteEdge: '',
            headerRewriteMid: '',
            headerRewriteRedirectRouter: '',
            notes: ''
        };
    };
    $scope.resetDSData();

    $scope.onDeliveryProtocolChange = function() {
        if ($scope.dsData.deliveryProtocol == 'http-to-https') {
            $scope.dsData.routingType = 'http'; // routing type must be http for http-to-https protocol
        }
    };

    $scope.onRoutingTypeChange = function() {
        if ($scope.dsData.routingType == 'dns') {
            $scope.dsData.headerRewriteRedirectRouter = ''; // not relevant for dns
        }
    };

    $scope.onNegativeCachingCustomizationChange = function() {
        if ($scope.dsData.hasNegativeCachingCustomization == false) {
            $scope.dsData.negativeCachingCustomizationNote = ''; // note no relevant
        }
    };

    $scope.onRateLimitingChange = function() {
        if ($scope.dsData.rateLimitingGBPS <= 0 && $scope.dsData.rateLimitingTPS <= 0) {
            $scope.dsData.overflowService = ''; // overflow service is irrelevant if no rate limits
        }
    };

    $scope.navigateToDashboard = function() {
        $location.url('/dashboard');
    };

    $scope.incomplete = function(forms) {
        var incomplete = false;
        if (!$scope.serviceDescStarted || !$scope.trafficStarted || !$scope.originStarted || !$scope.coreStarted || !$scope.serviceLimitStarted || !$scope.notesStarted || !$scope.confirmStarted) {
            return true;
        }
        _.each(forms, function(form) {
            if (form.$invalid) {
                incomplete = true;
            }
        });
        return incomplete;
    };

    $scope.jumpToTab = function(tab, startedFlag) {
        $scope.openTab = tab;
        $scope[startedFlag] = true;
    };

    $scope.submitRequest = function(dsData) {
        deliveryServiceService.createDSRequest(dsData).finally(
            function() { $anchorScroll(); } // scroll to top
        );
    };

    $scope.addAlias = function() {
        $scope.dsData.serviceAliases.push('');
    };

    $scope.removeAlias = function(index) {
        if (index > 0) { // no removing the first one
            $scope.dsData.serviceAliases.splice(index, 1);
        }
    };

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

DeliveryServiceNewController.$inject = ['$scope', '$location', '$anchorScroll', 'formUtils', 'deliveryServiceService'];
module.exports = DeliveryServiceNewController;
