var FormDeliveryServiceController = function(deliveryService, $scope, deliveryServiceService) {

    $scope.deliveryService = deliveryService;

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormDeliveryServiceController.$inject = ['deliveryService', '$scope', 'deliveryServiceService'];
module.exports = FormDeliveryServiceController;