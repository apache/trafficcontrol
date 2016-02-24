var UsersController = function(users, $scope) {

    $scope.users = users;

    angular.element(document).ready(function () {
        $('#usersTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

UsersController.$inject = ['users', '$scope'];
module.exports = UsersController;