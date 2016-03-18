var UserEditController = function($scope) {

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Update'
    };

};

UserEditController.$inject = ['$scope'];
module.exports = UserEditController;
