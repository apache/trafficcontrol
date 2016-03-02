var DeliveryServiceService = function($http, $q, ENV) {

    this.getDeliveryServices = function(ignoreLoadingBar) {
        var deferred = $q.defer();

        $http.get(ENV.api['base_url'] + "deliveryservice", { ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                deferred.resolve(result.response);
            })
            .error(function(fault) {
                deferred.reject(fault);
            });

        return deferred.promise;
    };


    this.getDeliveryService = function(dsId, ignoreLoadingBar) {
        var promise = $http.get(ENV.api['base_url'] + "deliveryservice/" + dsId, { ignoreLoadingBar: ignoreLoadingBar })
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