var UsersController = function(users, $scope) {

    $scope.users = users;

    $scope.predicate = 'fullName';
    $scope.reverse = false;

    $scope.query = {
        text: ''
    };

    // pagination
    $scope.currentPage = 1;
    $scope.usersPerPage = $scope.users.length;

    $scope.show = function(count) {
        $scope.usersPerPage = count;
    };

    $scope.search = function(user) {
        var query = $scope.query.text.toLowerCase(),
            fullName = (user.fullName) ? user.fullName.toLowerCase() : '',
            username = (user.username) ? user.username.toLowerCase() : '',
            email = (user.email) ? user.email.toLowerCase() : '',
            rolename = (user.rolename) ? user.rolename.toLowerCase() : '',
            isSubstring =
                (fullName.indexOf(query) !== -1) ||
                (username.indexOf(query) !== -1) ||
                (email.indexOf(query) !== -1) ||
                (rolename.indexOf(query) !== -1);

        return isSubstring;
    };

};

UsersController.$inject = ['users', '$scope'];
module.exports = UsersController;