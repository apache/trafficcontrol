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

var StatsService = function($http, $q, messageModel, ENV) {

    var edgeBandwidthRequest,
        edgeBandwidthSummaryRequest;

    var edgeTransactionsRequest,
        edgeTransactionsSummaryRequest;

    var retentionPeriodInDays; // used as an influxdb override, leave null if no override

    var displayTimoutError = function(options) {
        var msg = (angular.isDefined(options.message)) ? options.message : 'Request timeout. Please narrow your chart window.';
        if (options.status.toString().match(/^5\d[24]$/)) {
            // 502 or 504
            messageModel.setMessages([ { level: 'error', text: msg } ], false);
        }
    };

    this.getEdgeBandwidthBatch = function(deliveryService, start, end, interval, exclude, ignoreLoadingBar, showError) {
        var deferred = $q.defer();

        var url = ENV.apiEndpoint['1.2'] + "deliveryservice_stats.json",
            params = { deliveryServiceName: deliveryService.xmlId, metricType: 'kbps', serverType: 'edge', startDate: start.seconds(00).format(), endDate: end.seconds(00).format(), interval: interval, exclude: exclude, retentionPeriodInDays: retentionPeriodInDays };

        $http.get(url, { params: params, ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                deferred.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showError) displayTimoutError({ status: status });
                deferred.reject();
            });

        return deferred.promise;
    };

    this.getEdgeBandwidth = function(deliveryService, start, end, interval, exclude, ignoreLoadingBar, showError) {
        if (edgeBandwidthRequest) {
            edgeBandwidthRequest.reject();
        }
        edgeBandwidthRequest = $q.defer();

        var url = ENV.apiEndpoint['1.2'] + "deliveryservice_stats.json",
            params = { deliveryServiceName: deliveryService.xmlId, metricType: 'kbps', serverType: 'edge', startDate: start.seconds(00).format(), endDate: end.seconds(00).format(), interval: interval, exclude: exclude, retentionPeriodInDays: retentionPeriodInDays};

        $http.get(url, { params: params, timeout: edgeBandwidthRequest.promise, ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                edgeBandwidthRequest.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showError) displayTimoutError({ status: status });
                edgeBandwidthRequest.reject();
            });

        return edgeBandwidthRequest.promise;
    };

    this.getEdgeBandwidthSummary = function(deliveryService, start, end, interval, exclude, ignoreLoadingBar, showError) {
        if (edgeBandwidthSummaryRequest) {
            edgeBandwidthSummaryRequest.reject();
        }
        edgeBandwidthSummaryRequest = $q.defer();

        var url = ENV.apiEndpoint['1.2'] + "deliveryservice_stats.json",
            params = { deliveryServiceName: deliveryService.xmlId, metricType: 'kbps', serverType: 'edge', startDate: start.seconds(00).format(), endDate: end.seconds(00).format(), interval: interval, exclude: exclude, retentionPeriodInDays: retentionPeriodInDays };

        $http.get(url, { params: params, timeout: edgeBandwidthSummaryRequest.promise, ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                edgeBandwidthSummaryRequest.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showError) displayTimoutError({ status: status });
                edgeBandwidthSummaryRequest.reject();
            });

        return edgeBandwidthSummaryRequest.promise;
    };

    this.getEdgeTransactions = function(deliveryService, start, end, interval, exclude, ignoreLoadingBar, showError) {
        if (edgeTransactionsRequest) {
            edgeTransactionsRequest.reject();
        }
        edgeTransactionsRequest = $q.defer();

        var url = ENV.apiEndpoint['1.2'] + "deliveryservice_stats.json",
            params = { deliveryServiceName: deliveryService.xmlId, metricType: 'tps_total', serverType: 'edge', startDate: start.seconds(00).format(), endDate: end.seconds(00).format(), interval: interval, exclude: exclude, retentionPeriodInDays: retentionPeriodInDays };

        $http.get(url, { params: params, timeout: edgeTransactionsRequest.promise, ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                edgeTransactionsRequest.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showError) displayTimoutError({ status: status });
                edgeTransactionsRequest.reject();
            });

        return edgeTransactionsRequest.promise;
    };

    this.getEdgeTransactionsSummary = function(deliveryService, start, end, interval, ignoreLoadingBar, showError) {
        if (edgeTransactionsSummaryRequest) {
            edgeTransactionsSummaryRequest.reject();
        }
        edgeTransactionsSummaryRequest = $q.defer();

        var url = ENV.apiEndpoint['1.2'] + "deliveryservice_stats.json",
            params = { deliveryServiceName: deliveryService.xmlId, metricType: 'tps_total', serverType: 'edge', startDate: start.seconds(00).format(), endDate: end.seconds(00).format(), interval: interval, exclude: 'series', retentionPeriodInDays: retentionPeriodInDays };

        $http.get(url, { params: params, timeout: edgeTransactionsSummaryRequest.promise, ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                edgeTransactionsSummaryRequest.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showError) displayTimoutError({ status: status });
                edgeTransactionsSummaryRequest.reject();
            });

        return edgeTransactionsSummaryRequest.promise;
    };

    this.getEdgeTransactionsByStatusGroup = function(deliveryService, httpStatus, start, end, interval, exclude, ignoreLoadingBar, showError) {
        var request = $q.defer();

        var url = ENV.apiEndpoint['1.2'] + "deliveryservice_stats.json",
            params = { deliveryServiceName: deliveryService.xmlId, metricType: 'tps_' + httpStatus, serverType: 'edge', startDate: start.seconds(00).format(), endDate: end.seconds(00).format(), interval: interval, exclude: exclude, retentionPeriodInDays: retentionPeriodInDays };

        $http.get(url, { params: params, ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                request.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showError) displayTimoutError({ status: status });
                request.reject();
            });

        return request.promise;
    };

};

StatsService.$inject = ['$http', '$q', 'messageModel', 'ENV'];
module.exports = StatsService;