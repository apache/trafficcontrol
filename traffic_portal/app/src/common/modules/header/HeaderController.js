var HeaderController = function($scope, $state, $window, $anchorScroll, $uibModal, authService, portalService, userModel, deliveryServicesModel, propertiesModel) {

    $scope.isCollapsed = true;

    $scope.user = angular.copy(userModel.user);

    $scope.deliveryServices = deliveryServicesModel.deliveryServices;

    $scope.properties = propertiesModel.properties;

    $scope.isState = function(state) {
        return $state.current.name.indexOf(state) !== -1;
    };

    $scope.showDS = function(dsId) {
        $state.go('trafficPortal.private.deliveryService.view.overview.detail', { deliveryServiceId: dsId } );
    };

    $scope.releaseVersion = function() {
        portalService.getReleaseVersionInfo()
            .then(function(result) {
                $uibModal.open({
                    templateUrl: 'common/modules/release/version/release.version.tpl.html',
                    controller: 'ReleaseVersionController',
                    size: 'sm',
                    resolve: {
                        params: function () {
                            return result;
                        }
                    }
                });
            });
    };

    $scope.about = function() {
        if ($scope.properties.about.iframe) {
            $scope.navigateToState('trafficPortal.public.about', false);
        } else {
            $window.open(
                $scope.properties.about.url,
                '_blank'
            );
        }
    };

    $scope.navigateToState = function(to, reload) {
        $state.go(to, {}, { reload: reload });
    };

    $scope.logout = function() {
        authService.logout();
    };

    var scrollToTop = function() {
        $anchorScroll(); // hacky?
    };

    $scope.$on('userModel::userUpdated', function(event) {
        $scope.user = angular.copy(userModel.user);
    });

    var init = function () {
        scrollToTop();
    };
    init();
};

HeaderController.$inject = ['$scope', '$state', '$window', '$anchorScroll', '$uibModal', 'authService', 'portalService', 'userModel', 'deliveryServicesModel', 'propertiesModel'];
module.exports = HeaderController;
