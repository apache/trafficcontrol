var FormNewDeliveryServiceController = function(deliveryService, $scope, $controller, deliveryServiceService) {

    // extends the FormDeliveryServiceController to inherit common methods
    angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: deliveryService, $scope: $scope }));

    $scope.deliveryServiceName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(deliveryService) {
        deliveryServiceService.createDeliveryService(deliveryService);
    };

};

FormNewDeliveryServiceController.$inject = ['deliveryService', '$scope', '$controller', 'deliveryServiceService'];
module.exports = FormNewDeliveryServiceController;