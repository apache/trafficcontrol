var DeliveryServicesController = function($scope, deliveryServicesModel) {

    $scope.deliveryServices = deliveryServicesModel.deliveryServices;

    angular.element(document).ready(function () {
        $('#deliveryServicesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

DeliveryServicesController.$inject = ['$scope', 'deliveryServicesModel'];
module.exports = DeliveryServicesController;