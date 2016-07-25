var DeliveryServiceViewController = function($scope, deliveryService) {

    $scope.deliveryService = deliveryService;

};

DeliveryServiceViewController.$inject = ['$scope', 'deliveryService'];
module.exports = DeliveryServiceViewController;
