var UsersController = function(users, $scope) {

    $scope.users = users;

};

UsersController.$inject = ['users', '$scope'];
module.exports = UsersController;