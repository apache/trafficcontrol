var UserService = function(Restangular, $http, $location, $q, authService, userModel, messageModel, ENV) {

    var service = this;

    this.getCurrentUser = function() {
        var token = $location.search().token,
            deferred = $q.defer();

        if (angular.isDefined(token)) {
            $location.search('token', null); // remove the token query param
            authService.tokenLogin(token)
                .then(
                    function(response) {
                        service.getCurrentUser();
                    }
                );
        } else {
            $http.get(ENV.api['root'] + "user/current.json")
                .success(function(result) {
                    userModel.setUser(result.response);
                    deferred.resolve(result.response);
                })
                .error(function(fault) {
                    deferred.reject(fault);
                });

            return deferred.promise;
        }
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
        return Restangular.all('users').getList();
    };

    this.getUser = function(id) {
        return Restangular.one("users", id).get();
    };

    this.createUser = function(user) {
        return Restangular.service('users').post(user)
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
        return Restangular.one("users", id).remove()
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

UserService.$inject = ['Restangular', '$http', '$location', '$q', 'authService', 'userModel', 'messageModel', 'ENV'];
module.exports = UserService;