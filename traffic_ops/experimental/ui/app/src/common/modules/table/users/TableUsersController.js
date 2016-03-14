var TableUsersController = function(users, $scope, $location) {

    $scope.users = users;

    $scope.editUser = function(id) {
        $location.path($location.path() + '/' + id);
    };

    angular.element(document).ready(function () {
        $('#usersTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableUsersController.$inject = ['users', '$scope', '$location'];
module.exports = TableUsersController;