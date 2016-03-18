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

    this.createUser = function(user) {
        return Restangular.service('tm_user').post(user)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'User created' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'User create failed' } ], false);
                }
            );
    };

    this.updateUser = function(user) {
        return user.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'User updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'User update failed' } ], false);
                }
            );
    };

    this.deleteUser = function(id) {
        return Restangular.one("tm_user", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'User deleted' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'User delete failed' } ], false);
                }
            );
    };

};

UserService.$inject = ['Restangular', 'userModel', 'messageModel'];
module.exports = UserService;