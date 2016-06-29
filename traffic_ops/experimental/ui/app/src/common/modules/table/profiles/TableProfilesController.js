var TableProfilesController = function(profiles, $scope, locationUtils) {

    $scope.profiles = profiles;

    $scope.editProfile = function(id) {
        locationUtils.navigateToPath('/admin/profiles/' + id + '/edit');
    };

    $scope.createProfile = function() {
        locationUtils.navigateToPath('/admin/profiles/new');
    };

    angular.element(document).ready(function () {
        $('#profilesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableProfilesController.$inject = ['profiles', '$scope', 'locationUtils'];
module.exports = TableProfilesController;