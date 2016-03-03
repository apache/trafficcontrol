var HttpService = function($http, $q) {

    this.get = function(resource) {
        var deferred = $q.defer();

        $http.get(resource)
            .success(function(result) {
                deferred.resolve(result);
            })
            .error(function(fault) {
                deferred.reject(fault);
            });

        return deferred.promise;
    };

    this.post = function(resource, payload) {
        var deferred = $q.defer();

        $http.post(resource, payload)
            .success(function(result) {
                deferred.resolve(result);
            })
            .error(function(fault) {
                deferred.reject(fault);
            });

        return deferred.promise;
    };

    this.put = function(resource, payload) {
        var deferred = $q.defer();

        $http.put(resource, payload)
            .success(function(result) {
                deferred.resolve(result.response);
            })
            .error(function(fault) {
                deferred.reject(fault);
            });

        return deferred.promise;
    };

    this.delete = function(resource) {
        var deferred = $q.defer();

        $http.delete(resource)
            .success(function(result) {
                deferred.resolve(result.response);
            })
            .error(function(fault) {
                deferred.reject(fault);
            });

        return deferred.promise;
    };

};

HttpService.$inject = ['$http', '$q'];
module.exports = HttpService;