var FormEditProfileController = function(profile, $scope, $controller, $uibModal, $anchorScroll, locationUtils, profileService) {

    // extends the FormProfileController to inherit common methods
    angular.extend(this, $controller('FormProfileController', { profile: profile, $scope: $scope }));

    var deleteProfile = function(profile) {
        profileService.deleteProfile(profile.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/profiles');
            });
    };

    $scope.profileName = angular.copy(profile.name);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.save = function(profile) {
        profileService.updateProfile(profile).
            then(function() {
                $scope.profileName = angular.copy(profile.name);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(profile) {
        var params = {
            title: 'Delete Profile: ' + profile.name,
            key: profile.name
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
            controller: 'DialogDeleteController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deleteProfile(profile);
        }, function () {
            // do nothing
        });
    };

};

FormEditProfileController.$inject = ['profile', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'profileService'];
module.exports = FormEditProfileController;