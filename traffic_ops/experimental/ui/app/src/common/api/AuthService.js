var AuthService = function($http, $state, $location, $q, deliveryServicesModel, userModel, messageModel, ENV) {

    this.login = function(username, password) {
        userModel.resetUser();
        deliveryServicesModel.resetDeliveryServices();
        var promise = $http.post(
                ENV.apiEndpoint['1.2'] + "user/login", { u: username, p: password })
            .success(function(result) {
                var redirect = decodeURIComponent($location.search().redirect);
                if (redirect !== 'undefined') {
                    $location.search('redirect', null); // remove the redirect query param
                    $location.url(redirect);
                } else {
                    $location.url('/dashboards/one');
                }
                return result;
            })
            .error(function(fault) {
                return fault;
            });

        return promise;
    };

    this.logout = function() {
        userModel.resetUser();
        deliveryServicesModel.resetDeliveryServices();
        var promise = $http.post(
                ENV.apiEndpoint['1.2'] + "user/logout")
            .success(function(result) {
                if ($state.current.name == 'trafficOps.public.login') {
                    messageModel.setMessages(result.alerts, false);
                } else {
                    messageModel.setMessages(result.alerts, true);
                    $state.go('trafficOps.public.login');
                }
                return result;
            })
            .error(function(fault) {
                return fault;
            });

        return promise;
    };

};

AuthService.$inject = ['$http', '$state', '$location', '$q', 'deliveryServicesModel', 'userModel', 'messageModel', 'ENV'];
module.exports = AuthService;