var FormDeliveryServiceController = function(deliveryService, $scope, formUtils, deliveryServiceService) {

    $scope.deliveryService = deliveryService;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormDeliveryServiceController.$inject = ['deliveryService', '$scope', 'formUtils', 'deliveryServiceService'];
module.exports = FormDeliveryServiceController;