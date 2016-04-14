var ProfileService = function(Restangular, messageModel) {

    this.getProfiles = function() {
        return Restangular.all('profile').getList();
    };

    this.getProfile = function(id) {
        return Restangular.one("profile", id).get();
    };

    this.createProfile = function(profile) {
        return Restangular.service('profile').post(profile)
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Profile created' } ], true);
            },
            function() {
                messageModel.setMessages([ { level: 'error', text: 'Profile create failed' } ], false);
            }
        );
    };

    this.updateProfile = function(profile) {
        return profile.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Profile updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Profile update failed' } ], false);
                }
        );
    };

    this.deleteProfile = function(id) {
        return Restangular.one("profile", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Profile deleted' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Profile delete failed' } ], false);
                }
        );
    };

};

ProfileService.$inject = ['Restangular', 'messageModel'];
module.exports = ProfileService;