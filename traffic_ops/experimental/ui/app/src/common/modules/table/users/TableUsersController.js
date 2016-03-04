var TableUsersController = function(users, $scope, $location) {

    $scope.users = users.response;

    $scope.editUser = function(id) {
        console.log('/administer/users/' + id);
        $location.url('/administer/users/' + id);
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