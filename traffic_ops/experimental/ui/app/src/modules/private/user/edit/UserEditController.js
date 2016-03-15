var UserEditController = function(showDelete, $scope) {

    $scope.showDelete = showDelete;

};

UserEditController.$inject = ['showDelete', '$scope'];
module.exports = UserEditController;
