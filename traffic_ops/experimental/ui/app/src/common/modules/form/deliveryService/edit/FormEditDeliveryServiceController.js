var FormEditDeliveryServiceController = function(deliveryService, $scope, $controller, $uibModal, $anchorScroll, locationUtils, deliveryServiceService) {

    // extends the FormDeliveryServiceController to inherit common methods
    angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: deliveryService, $scope: $scope }));

    var deleteDeliveryService = function(deliveryService) {
        deliveryServiceService.deleteDeliveryService(deliveryService.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/delivery-services');
            });
    };

    $scope.deliveryServiceName = angular.copy(deliveryService.displayName);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(deliveryService) {
        deliveryServiceService.updateDeliveryService(deliveryService).
            then(function() {
                $scope.deliveryServiceName = angular.copy(deliveryService.displayName);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(deliveryService) {
        var params = {
            title: 'Delete Delivery Service: ' + deliveryService.displayName,
            key: deliveryService.xmlId
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
            controller: 'DialogDeleteController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deleteDeliveryService(deliveryService);
        }, function () {
            // do nothing
        });
    };

};

FormEditDeliveryServiceController.$inject = ['deliveryService', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'deliveryServiceService'];
module.exports = FormEditDeliveryServiceController;