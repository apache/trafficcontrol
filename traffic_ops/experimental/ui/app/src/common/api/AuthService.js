var AuthService = function($http, $state, $location, $q, $state, httpService, userModel, messageModel, ENV) {

    this.login = function(username, password) {
        userModel.resetUser();
        return httpService.post(ENV.api['root'] + 'user/login', { u: username, p: password })
            .then(
                function(result) {
                    var redirect = decodeURIComponent($location.search().redirect);
                    if (redirect !== 'undefined') {
                        $location.search('redirect', null); // remove the redirect query param
                        $location.url(redirect);
                    } else {
                        $location.url('/monitor/dashboards/one');
                    }
                },
                function(fault) {
                    // do nothing
                }
            );
    };

    this.tokenLogin = function(token) {
        userModel.resetUser();
        return httpService.post(ENV.api['root'] + 'user/login/token', { t: token });
    };

    this.logout = function() {
        userModel.resetUser();
        httpService.post(ENV.api['root'] + 'user/logout').
            then(
                function(result) {
                    if ($state.current.name == 'trafficOps.public.login') {
                        messageModel.setMessages(result.alerts, false);
                    } else {
                        messageModel.setMessages(result.alerts, true);
                        $state.go('trafficOps.public.login');
                    }
                    return result;
                }
        );
    };

    this.resetPassword = function(email) {
        // Todo: api endpoint not implemented yet
    };

};

AuthService.$inject = ['$http', '$state', '$location', '$q', '$state', 'httpService', 'userModel', 'messageModel', 'ENV'];
module.exports = AuthService;