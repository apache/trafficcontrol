var DeliveryServiceService = function($http, $q, ENV) {

    this.getDeliveryServices = function(ignoreLoadingBar) {
        var deferred = $q.defer();

        $http.get(ENV.apiEndpoint['1.2'] + "deliveryservices.json", { ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                deferred.resolve(result.response);
            })
            .error(function(fault) {
                deferred.reject(fault);
            });

        return deferred.promise;
    };


    this.getDeliveryService = function(dsId, ignoreLoadingBar) {
        var promise = $http.get(ENV.apiEndpoint['1.2'] + "deliveryservices/" + dsId + ".json", { ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                return result;
            })
            .error(function(fault) {
            });

        return promise;
    };


};

DeliveryServiceService.$inject = ['$http', '$q', 'ENV'];
module.exports = DeliveryServiceService;