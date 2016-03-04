var AuthService = function($http, $state, $location, $q, jwtHelper, httpService, userModel, messageModel, ENV) {

    this.login = function(username, password) {
        userModel.resetUser();
        return httpService.post(ENV.api['root'] + 'login', { u: username, p: password })
            .then(
                function(result) {
                    userModel.setToken(result.Token);
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
        $state.go('trafficOps.public.login');
        // Todo: api endpoint not implemented yet
    };

    this.resetPassword = function(email) {
        // Todo: api endpoint not implemented yet
    };

};

AuthService.$inject = ['$http', '$state', '$location', '$q', 'jwtHelper', 'httpService', 'userModel', 'messageModel', 'ENV'];
module.exports = AuthService;