/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 */

var AuthService = function($http, $state, $location, $q, userModel, deliveryServicesModel, messageModel, ENV) {

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
                    $location.url('/dashboard');
                }
                return result;
            })
            .error(function(fault) {
                return fault;
            });

        return promise;
    };

    this.tokenLogin = function(token) {
        userModel.resetUser();
        deliveryServicesModel.resetDeliveryServices();
        var deferred = $q.defer();
        $http.post(
                ENV.apiEndpoint['1.2'] + "user/login/token", { t: token })
            .success(function(result) {
                deferred.resolve(result);
            })
            .error(function(fault) {
                deferred.reject(fault);
            });

        return deferred.promise;
    };

    this.logout = function() {
        userModel.resetUser();
        deliveryServicesModel.resetDeliveryServices();
        var promise = $http.post(
                ENV.apiEndpoint['1.2'] + "user/logout")
            .success(function(result) {
                if ($state.current.name == 'trafficPortal.public.home.landing') {
                    messageModel.setMessages(result.alerts, false);
                } else {
                    messageModel.setMessages(result.alerts, true);
                    $state.go('trafficPortal.public.home.landing');
                }
                return result;
            })
            .error(function(fault) {
                return fault;
            });

        return promise;
    };

};

AuthService.$inject = ['$http', '$state', '$location', '$q', 'userModel', 'deliveryServicesModel', 'messageModel', 'ENV'];
module.exports = AuthService;