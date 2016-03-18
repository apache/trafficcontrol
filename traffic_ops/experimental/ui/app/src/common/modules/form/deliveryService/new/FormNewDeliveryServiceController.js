var FormNewDeliveryServiceController = function(deliveryService, $scope, $controller, locationUtils, deliveryServiceService) {

    // extends the FormDeliveryServiceController to inherit common methods
    angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: deliveryService, $scope: $scope }));

    $scope.deliveryServiceName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(deliveryService) {
        deliveryServiceService.createDeliveryService(deliveryService).
            then(function() {
                locationUtils.navigateToPath('/configure/delivery-services');
            });
    };

};

FormNewDeliveryServiceController.$inject = ['deliveryService', '$scope', '$controller', 'locationUtils', 'deliveryServiceService'];
module.exports = FormNewDeliveryServiceController;