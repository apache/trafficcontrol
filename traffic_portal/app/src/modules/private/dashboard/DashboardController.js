var DashboardController = function($scope, $location, chartModel, deliveryServicesModel, propertiesModel, userModel) {

    $scope.deliveryServices = deliveryServicesModel.deliveryServices;

    $scope.properties = propertiesModel.properties;

    $scope.user = angular.copy(userModel.user);

    $scope.requestDS = function() {
        $location.url('/delivery-service/new');
    };

    var init = function () {
        chartModel.resetChart(); // set chart back to default parameters
    };
    init();
};

DashboardController.$inject = ['$scope', '$location', 'chartModel', 'deliveryServicesModel', 'propertiesModel', 'userModel'];
module.exports = DashboardController;
