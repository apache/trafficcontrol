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

var UserService = function($http, $state, $q, $location, authService, userModel, messageModel, ENV) {

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
            $http.get(ENV.apiEndpoint['1.2'] + "user/current.json")
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

    this.updateCurrentUser = function(userData) {
        var deferred = $q.defer();
        var user = _.omit(userData, 'loaded'); // let's pull the loaded key off of there
        $http.post(ENV.apiEndpoint['1.2'] + "user/current/update", { user: user })
            .success(function(result) {
                userModel.setUser(userData);
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

    this.createUserPurgeJob = function(jobParams) {
        var deferred = $q.defer();
        $http.post(ENV.apiEndpoint['1.2'] + "user/current/jobs", jobParams)
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