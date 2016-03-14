var TableDeliveryServicesController = function(deliveryServices, $scope, $location) {

    $scope.deliveryServices = deliveryServices;

    $scope.editDeliveryService = function(id) {
        $location.path($location.path() + '/' + id);
    };

    angular.element(document).ready(function () {
        $('#deliveryServicesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableDeliveryServicesController.$inject = ['deliveryServices', '$scope', '$location'];
module.exports = TableDeliveryServicesController;