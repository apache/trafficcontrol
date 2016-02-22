var ConfigureController = function($scope, $location) {

    $scope.navigateToPath = function(path) {
        $location.url(path);
    };

};

ConfigureController.$inject = ['$scope', '$location'];
module.exports = ConfigureController;