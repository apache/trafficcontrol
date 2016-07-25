var DeliveryServiceController = function($scope, userModel) {

    $scope.user = angular.copy(userModel.user);

};

DeliveryServiceController.$inject = ['$scope', 'userModel'];
module.exports = DeliveryServiceController;
