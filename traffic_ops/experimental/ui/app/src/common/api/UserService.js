var UserService = function($http, $state, $q, $location, authService, userModel, messageModel, ENV) {

    this.getCurrentUser = function() {
        var deferred = $q.defer();

        $http.get(ENV.apiEndpoint['1.2'] + "user/current.json")
            .success(function(result) {
                userModel.setUser(result.response);
                deferred.resolve(result.response);
            })
            .error(function(fault) {
                deferred.reject(fault);
            });

        return deferred.promise;
    };

    this.updateCurrentUser = function(user) {
        var deferred = $q.defer();
        $http.post(ENV.apiEndpoint['1.2'] + "user/current/update", { user: user })
            .success(function(result) {
                userModel.setUser(user);
                messageModel.setMessages(result.alerts, false);
                deferred.resolve(result);
            })
            .error(function(fault) {
                if (angular.isDefined(fault.alerts)) {
                    messageModel.setMessages(fault.alerts, false);
                }
                deferred.reject();
            });

        return deferred.promise;
    };

    this.resetPassword = function(email) {
        var deferred = $q.defer();
        $http.post(
                ENV.apiEndpoint['1.2'] + "user/reset_password", { email: email })
            .success(function(result) {
                messageModel.setMessages(result.alerts, false);
                deferred.resolve(result);
            })
            .error(function(fault) {
                if (angular.isDefined(fault.alerts)) {
                    messageModel.setMessages(fault.alerts, false);
                }
                deferred.reject();
            });

        return deferred.promise;
    };

};

UserService.$inject = ['$http', '$state', '$q', '$location', 'authService', 'userModel', 'messageModel', 'ENV'];
module.exports = UserService;