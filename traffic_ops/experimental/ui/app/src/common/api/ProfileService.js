var ProfileService = function(Restangular, locationUtils, messageModel) {

    this.getProfiles = function() {
        return Restangular.all('profiles').getList();
    };

    this.getProfile = function(id) {
        return Restangular.one("profiles", id).get();
    };

    this.createProfile = function(profile) {
        return Restangular.service('profiles').post(profile)
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Profile created' } ], true);
                locationUtils.navigateToPath('/admin/profiles');
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.updateProfile = function(profile) {
        return profile.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Profile updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
        );
    };

    this.deleteProfile = function(id) {
        return Restangular.one("profiles", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Profile deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
        );
    };

};

ProfileService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = ProfileService;