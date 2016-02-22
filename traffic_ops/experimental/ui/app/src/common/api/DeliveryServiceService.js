var DeliveryServiceService = function($http, $log, deliveryServicesModel, ENV) {

    this.getDeliveryServices = function(ignoreLoadingBar) {
        var promise = $http.get(ENV.apiEndpoint['1.2'] + "deliveryservices.json", { ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                deliveryServicesModel.setDeliveryServices(result.response);
                return result.response;
            })
            .error(function(fault) {
            });

        return promise;
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

DeliveryServiceService.$inject = ['$http', '$log', 'deliveryServicesModel', 'ENV'];
module.exports = DeliveryServiceService;