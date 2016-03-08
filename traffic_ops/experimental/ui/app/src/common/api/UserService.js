var UserService = function(Restangular, userModel, messageModel) {

    this.getCurrentUser = function() {
        return Restangular.one("tm_user", userModel.userId).get()
            .then(function(user) {
                userModel.setUser(user);
            });
    };

    this.updateCurrentUser = function(user) {
        return user.put()
            .then(
                function() {
                    userModel.setUser(user);
                    messageModel.setMessages([ { level: 'success', text: 'User updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'User updated failed' } ], false);
                }
            );
    };

    this.getUsers = function() {
        return Restangular.all('tm_user').getList();
    };

    this.getUser = function(id) {
        return Restangular.one("tm_user", id).get();
    };

    this.updateUser = function(user) {
        return user.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'User updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'User updated failed' } ], false);
                }
            );
    };

};

UserService.$inject = ['Restangular', 'userModel', 'messageModel'];
module.exports = UserService;