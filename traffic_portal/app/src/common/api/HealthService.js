var HealthService = function($http, $q, ENV) {

    var displayTimoutError = function(options) {
        var msg = (angular.isDefined(options.message)) ? options.message : 'Request timeout. Failed to load cache groups.';
        if (options.status.toString().match(/^5\d[24]$/)) {
            // 502 or 504
            messageModel.setMessages([ { level: 'error', text: msg } ], false);
        }
    };

    this.getCacheGroupHealth = function(entityId, ignoreLoadingBar, showTimeoutError) {
        var deferred = $q.defer();
        $http.get(ENV.apiEndpoint['1.2'] + "cdns/health.json", { ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                deferred.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showTimeoutError) displayTimoutError({ status: status });
                deferred.reject();
            });

        return deferred.promise;
    };

};

HealthService.$inject = ['$http', '$q', 'ENV'];
module.exports = HealthService;