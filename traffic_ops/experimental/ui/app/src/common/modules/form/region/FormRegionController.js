var FormRegionController = function(region, $scope, regionService) {

    $scope.region = region;

    $scope.hasError = function(input) {
        return !input.$focused && input.$dirty && input.$invalid;
    };

    $scope.hasPropertyError = function(input, property) {
        return !input.$focused && input.$dirty && input.$error[property];
    };

};

FormRegionController.$inject = ['region', '$scope', 'regionService'];
module.exports = FormRegionController;