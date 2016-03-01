var AuthService = function($http, $state, $location, $q, httpService, userModel, messageModel, ENV) {

    this.login = function(username, password) {
        userModel.resetUser();
        return httpService.post(ENV.apiEndpoint['login'], { u: username, p: password })
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

    this.logout = function() {
        userModel.resetUser();
        return httpService.post(ENV.apiEndpoint['logout'])
            .then(
                function(result) {
                    if ($state.current.name == 'trafficOps.public.login') {
                        messageModel.setMessages(result.alerts, false);
                    } else {
                        messageModel.setMessages(result.alerts, true);
                        $state.go('trafficOps.public.login');
                    }
                },
                function(fault) {
                    // do nothing
                }
            );
    };

    this.resetPassword = function(email) {
        return httpService.post(ENV.apiEndpoint['reset_password'], { email: email })
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

AuthService.$inject = ['$http', '$state', '$location', '$q', 'httpService', 'userModel', 'messageModel', 'ENV'];
module.exports = AuthService;