var FormNewProfileController = function(profile, $scope, $controller, profileService) {

    // extends the FormProfileController to inherit common methods
    angular.extend(this, $controller('FormProfileController', { profile: profile, $scope: $scope }));

    $scope.profileName = 'New';

    $scope.settings = {
        showDelete: false,
        saveLabel: 'Create'
    };

    $scope.save = function(profile) {
        profileService.createProfile(profile);
    };

};

FormNewProfileController.$inject = ['profile', '$scope', '$controller', 'profileService'];
module.exports = FormNewProfileController;