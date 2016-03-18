var TableDeliveryServicesController = function(deliveryServices, $scope, locationUtils) {

    $scope.deliveryServices = deliveryServices;

    $scope.editDeliveryService = function(id) {
        locationUtils.navigateToPath('/configure/delivery-services/' + id + '/edit');
    };

    $scope.createDeliveryService = function() {
        locationUtils.navigateToPath('/configure/delivery-services/new');
    };

    angular.element(document).ready(function () {
        $('#deliveryServicesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableDeliveryServicesController.$inject = ['deliveryServices', '$scope', 'locationUtils'];
module.exports = TableDeliveryServicesController;