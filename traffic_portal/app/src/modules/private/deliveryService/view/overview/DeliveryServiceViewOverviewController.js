var DeliveryServiceViewOverviewController = function($scope, $location, $state, $uibModal, propertiesModel, deliveryServiceService, chartModel) {

    var getFailoverStatus = function() {
        var ignoreLoadingBar = true;
        deliveryServiceService.getState($scope.deliveryService.id, ignoreLoadingBar)
            .then(function(response) {
                $scope.failover = response.failover;
            });
    };

    $scope.properties = propertiesModel.properties;

    $scope.failover = {
        configured: false,
        enabled: false,
        destination: {
            location: null,
            type: ''
        },
        locations: []
    };

    $scope.viewConfig = function(ds) {

        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/deliveryService/config/edit/deliveryService.config.edit.tpl.html',
            controller: 'DSConfigEditController',
            size: 'lg',
            windowClass: 'ds-config-modal',
            resolve: {
                deliveryService: function () {
                    return angular.copy(ds);
                }
            }
        });

        modalInstance.result.then(function() {
        }, function () {
            // do nothing
        });
    };

    $scope.navigateToChart = function(dsId, type) {
        $location.url('/delivery-service/' + dsId + '/chart/' + type).search({ start: moment(chartModel.chart.start).format(), end: moment(chartModel.chart.end).format() });
    };

    angular.element(document).ready(function () {
        if ($scope.deliveryService && $scope.deliveryService.active) {
            getFailoverStatus();
        }
    });

};

DeliveryServiceViewOverviewController.$inject = ['$scope', '$location', '$state', '$uibModal', 'propertiesModel', 'deliveryServiceService', 'chartModel'];
module.exports = DeliveryServiceViewOverviewController;
