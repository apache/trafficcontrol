var TableUsersController = function(users, $scope, locationUtils) {

    $scope.users = users;

    $scope.editUser = function(id) {
        locationUtils.navigateToPath('/admin/users/' + id + '/edit');
    };

    $scope.createUser = function() {
        locationUtils.navigateToPath('/admin/users/new');
    };

    angular.element(document).ready(function () {
        $('#usersTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableUsersController.$inject = ['users', '$scope', 'locationUtils'];
module.exports = TableUsersController;