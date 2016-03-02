var TableDeliveryServicesController = function(deliveryServices, $scope) {

    $scope.deliveryServices = deliveryServices.response;

    angular.element(document).ready(function () {
        $('#deliveryServicesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableDeliveryServicesController.$inject = ['deliveryServices', '$scope'];
module.exports = TableDeliveryServicesController;