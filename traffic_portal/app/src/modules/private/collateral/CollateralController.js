var CollateralController = function($scope, propertiesModel) {

    $scope.collateralItems = propertiesModel.properties.collateral.items;

};

CollateralController.$inject = ['$scope', 'propertiesModel'];
module.exports = CollateralController;
