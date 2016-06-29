var FormNewProfileController = function(profile, $scope, $controller, locationUtils, profileService) {

    // extends the FormProfileController to inherit common methods
    angular.extend(this, $controller('FormProfileController', { profile: profile, $scope: $scope }));

    $scope.profileName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(profile) {
        profileService.createProfile(profile).
            then(function() {
                locationUtils.navigateToPath('/admin/profiles');
            });
    };

};

FormNewProfileController.$inject = ['profile', '$scope', '$controller', 'locationUtils', 'profileService'];
module.exports = FormNewProfileController;