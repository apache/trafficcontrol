var ConfigureDeliveryServicesController = function($scope, $interval, deliveryServiceService, deliveryServicesModel) {

    var refreshInterval;

    var refreshDeliveryServices = function() {
        deliveryServiceService.getDeliveryServices(true);
    };

    $scope.deliveryServicesModel = deliveryServicesModel;

    $scope.predicate = 'xmlId';
    $scope.reverse = false;

    $scope.query = {
        text: ''
    };

    // pagination
    $scope.currentPage = 1;
    $scope.dsPerPage = $scope.deliveryServicesModel.deliveryServices.length;

    $scope.show = function(count) {
        $scope.dsPerPage = count;
    };

    $scope.search = function(ds) {
        var query = $scope.query.text.toLowerCase(),
            xmlId = ds.xmlId.toLowerCase(),
            orgServerFqdn = ds.orgServerFqdn.toLowerCase(),
            isSubstring = (xmlId.indexOf(query) !== -1) || (orgServerFqdn.indexOf(query) !== -1);

        return isSubstring;
    };

    angular.element(document).ready(function () {
        refreshInterval = $interval(function() { refreshDeliveryServices() }, 1 * 60 * 1000); // every 1 min delivery services will refresh
    });

    $scope.$on("$destroy", function() {
        if (angular.isDefined(refreshInterval)) {
            $interval.cancel(refreshInterval);
            refreshInterval = undefined;
        }
    });

};

ConfigureDeliveryServicesController.$inject = ['$scope', '$interval', 'deliveryServiceService', 'deliveryServicesModel'];
module.exports = ConfigureDeliveryServicesController;