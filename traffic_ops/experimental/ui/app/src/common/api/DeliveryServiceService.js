var DeliveryServiceService = function($http, $q, httpService, ENV) {

    this.getDeliveryServices = function(endpoint) {
        return httpService.get(endpoint);
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

DeliveryServiceService.$inject = ['$http', '$q', 'httpService', 'ENV'];
module.exports = DeliveryServiceService;