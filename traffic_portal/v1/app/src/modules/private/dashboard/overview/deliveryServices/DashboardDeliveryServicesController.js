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

var DashboardDeliveryServicesController = function($window, $rootScope, $scope, $location, $uibModal, $q, $interval, $filter, $anchorScroll, dateUtils, numberUtils, chartUtils, messageModel, chartModel, statsService) {

    var bandwidthInterval;

    var createDeliveryServiceSparkline = function(dsId, unitSize) {
        var $sparkLine = $(".delivery-service-sparkline-" + dsId),
            data = chartUtils.formatData($sparkLine.data('sparkline'));

        var options = {
            xaxis: {
                mode: "time",
                timezone: "browser",
                ticks: false
            },
            yaxis: { ticks: false },
            grid: {
                borderWidth: 0,
                hoverable: true
            },
            tooltip: {
                show: true,
                content: function(label, xval, yval, flotItem){
                    return '<span>' + $filter('number')(yval, 2) + ' ' + unitSize + 'ps @ ' + dateUtils.dateFormat(xval, "h:MM:ss tt (Z)") + '</span>';
                }
            }
        };

        var series = [
            {
                data: data,
                color: '#337ab7',
                lines: {
                    lineWidth: 0.8
                },
                shadowSize: 0
            }
        ];

        // draw the sparkline
        $.plot($sparkLine, series, options);
    };

    $scope.loadBandwidth = function(dsId, showLoading) {
        if (!$rootScope.online) return;

        var deliveryService = _.find($scope.deliveryServices, function(ds){ return ds.id === dsId }),
            $sparkLine = $(".delivery-service-sparkline-" + dsId),
            yesterday = moment().subtract(24, 'hours'),
            now = moment(),
            exclude = '';

        if (showLoading) {
            try {
                $sparkLine.empty(); // remove the sparkline if it's there
            } catch (err) {
                // there was no sparkline evidently so no need to clear it
            }

            $sparkLine.html("<div>Loading...</div>");
            $('.delivery-service-last-' + dsId).html("<div>Calculating...</div>");
        }

        statsService.getEdgeBandwidthBatch(deliveryService, yesterday, now, '60s', exclude, true, false)
            .then(
                function(response) {

                    var sparklineData = [],
                        originalValue = 0,
                        convertedValue = 0,
                        summary = response.summary,
                        series = response.series;

                    try {
                        $scope.unitSize = numberUtils.shrink(summary.average)[1];
                        _.each(series.values, function(seriesItem) {
                            if (_.isNumber(seriesItem[1])) {
                                originalValue = seriesItem[1];
                                convertedValue = numberUtils.convertTo(seriesItem[1], $scope.unitSize);
                                sparklineData.push(moment(seriesItem[0]).valueOf()); // time in milliseconds
                                sparklineData.push(convertedValue); // value
                            }
                        });
                    }
                    catch (e) {
                        // no bandwidth for delivery service
                    }

                    deliveryService.last = originalValue;
                    var convertedLast = numberUtils.shrink(originalValue);
                    $('.delivery-service-last-' + dsId).html($filter('number')(convertedLast[0], 2) + ' ' + convertedLast[1] + 'ps');
                    $sparkLine.data('sparkline', sparklineData.join(','));
                    createDeliveryServiceSparkline(dsId, $scope.unitSize);
                },
                function(fault) {
                    $('.delivery-service-sparkline-' + dsId).html('Error');
                    $('.delivery-service-last-' + dsId).html('Error');
                });
    };

    $scope.search = {
        query: ""
    };

    $scope.predicate = 'displayName';
    $scope.reverse = false;

    $scope.dsOptions = {
        inactive: true
    };

    $scope.unitSize = 'Kb';

    $scope.autoLoadLimit = 10;

    $scope.navigateToDeliveryService = function(dsId) {
        $location.url('/delivery-service/' + dsId).search({ start: moment(chartModel.chart.start).format(), end: moment(chartModel.chart.end).format() });
    };

    $scope.viewConfig = function(ds) {

        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/deliveryService/config/edit/deliveryService.config.edit.tpl.html',
            controller: 'DSConfigEditController',
            size: 'lg',
            windowClass: 'ds-config-modal',
            resolve: {
                deliveryService: function (deliveryServiceService) {
                    return deliveryServiceService.getDeliveryService(ds.id);
                }
            }
        });

        modalInstance.result.then(function() {
        }, function () {
            // do nothing
        });
    };

    $scope.showInactive = function(show) {
        $scope.dsOptions.inactive = show;
    };

    $scope.hideDeliveryService = function(ds) {
        var query = $scope.search.query.toLowerCase(),
            id = ds.id.toString(),
            xmlId = ds.xmlId.toLowerCase(),
            displayName = ds.displayName.toLowerCase(),
            isSubstring = (id.indexOf(query) !== -1) || (xmlId.indexOf(query) !== -1) || (displayName.indexOf(query) !== -1);

        return !isSubstring || ($scope.dsOptions.inactive == false && !ds.active);
    };

    angular.element(document).ready(function () {
        // if you do not exceed # of delivery services allowed for autoload, we'll autoload bandwidth for each delivery service and refresh on a timer
        if ($scope.deliveryServices.length <= $scope.autoLoadLimit) {
            $scope.loadAllBandwidth();
            bandwidthInterval = $interval($scope.loadAllBandwidth, (5*60*1000)); // new bandwidth data every 5 minutes
        }
    });

    $scope.$on("$destroy", function() {
        if (angular.isDefined(bandwidthInterval)) {
            $interval.cancel(bandwidthInterval);
            bandwidthInterval = undefined;
        }
    });

};

DashboardDeliveryServicesController.$inject = ['$window', '$rootScope', '$scope', '$location', '$uibModal', '$q', '$interval', '$filter', '$anchorScroll', 'dateUtils', 'numberUtils', 'chartUtils', 'messageModel', 'chartModel', 'statsService'];
module.exports = DashboardDeliveryServicesController;
