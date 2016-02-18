var DeliveryServiceService = function($http, deliveryServicesModel, ENV) {

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

};

DeliveryServiceService.$inject = ['$http', 'deliveryServicesModel', 'ENV'];
module.exports = DeliveryServiceService;