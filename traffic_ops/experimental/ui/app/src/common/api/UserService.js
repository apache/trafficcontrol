var UserService = function($http, $state, $q, $location, authService, httpService, userModel, messageModel) {

    this.getCurrentUser = function(endpoint) {
        return httpService.get(endpoint)
            .then(function(result) {
                userModel.setUser(result.response[0]);
            });
    };

    this.updateCurrentUser = function(endpoint, user) {
        return httpService.put(endpoint, { user: user })
            .then(
                function(result) {
                    userModel.setUser(user);
                    messageModel.setMessages(result.alerts, false);
                },
                function(fault) {
                    messageModel.setMessages(fault.alerts, false);
                }
            );
    };

    this.getUsers = function(endpoint) {
        return httpService.get(endpoint);
    };


    this.getUser = function(endpoint) {
        return httpService.get(endpoint);
    };

    this.updateUser = function(endpoint, user) {
        return httpService.put(endpoint, { user: user })
            .then(
                function(result) {
                    messageModel.setMessages(result.alerts, false);
                },
                function(fault) {
                    messageModel.setMessages(fault.alerts, false);
                }
            );
    };

};

UserService.$inject = ['$http', '$state', '$q', '$location', 'authService', 'httpService', 'userModel', 'messageModel'];
module.exports = UserService;