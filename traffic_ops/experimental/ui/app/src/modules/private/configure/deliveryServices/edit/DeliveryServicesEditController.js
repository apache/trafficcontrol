var DeliveryServicesEditController = function(deliveryService, $scope) {

    $scope.deliveryService = deliveryService.data.response[0];

};

DeliveryServicesEditController.$inject = ['deliveryService', '$scope'];
module.exports = DeliveryServicesEditController;