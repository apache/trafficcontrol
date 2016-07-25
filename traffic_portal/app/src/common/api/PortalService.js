var PortalService = function($http, $q) {

    this.getReleaseVersionInfo = function() {
        var deferred = $q.defer();
        $http.get('traffic_portal_release.json')
            .success(function(result) {
                deferred.resolve(result);
            });

        return deferred.promise;
    };

    this.getProperties = function() {
        var deferred = $q.defer();
        $http.get('traffic_portal_properties.json')
            .success(function(result) {
                deferred.resolve(result.properties);
            });

        return deferred.promise;
    };

};

PortalService.$inject = ['$http', '$q'];
module.exports = PortalService;