var TrafficOpsService = function($http, $q) {

    this.getReleaseVersionInfo = function() {
        var deferred = $q.defer();
        $http.get('traffic_ops_release.json')
            .success(function(result) {
                deferred.resolve(result);
            })
            .error(function(fault) {
                deferred.reject(fault);
            });

        return deferred.promise;
    };

};

TrafficOpsService.$inject = ['$http', '$q'];
module.exports = TrafficOpsService;