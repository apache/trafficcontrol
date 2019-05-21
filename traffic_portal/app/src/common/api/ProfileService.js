/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

var ProfileService = function(Restangular, $http, $q, locationUtils, messageModel, ENV) {

    this.getProfiles = function(queryParams) {
        return Restangular.all('profiles').getList(queryParams);
    };

    this.getProfile = function(id, queryParams) {
        return Restangular.one("profiles", id).get(queryParams);
    };

    this.createProfile = function(profile) {
        return Restangular.service('profiles').post(profile)
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Profile created' } ], true);
                locationUtils.navigateToPath('/profiles');
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.updateProfile = function(profile) {
        return profile.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Profile updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
        );
    };

    this.deleteProfile = function(id) {
        var request = $q.defer();

        $http.delete(ENV.api['root'] + "profiles/" + id)
            .then(
                function(result) {
                    request.resolve(result.data);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                    request.reject(fault);
                }
            );

        return request.promise;
    };

    this.getParameterProfiles = function(paramId) {
        return Restangular.one('parameters', paramId).getList('profiles');
    };

    this.getParamUnassignedProfiles = function(paramId) {
        return Restangular.one('parameters', paramId).getList('unassigned_profiles');
    };

    this.cloneProfile = function(sourceName, cloneName) {
        return $http.post(ENV.api['root'] + "profiles/name/" + cloneName + "/copy/" + sourceName)
            .then(
                function(result) {
                    messageModel.setMessages(result.data.alerts, true);
                    locationUtils.navigateToPath('/profiles/' + result.data.response.id);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.exportProfile = function(id) {
        var deferred = $q.defer();

        $http.get(ENV.api['root'] + "profiles/" + id + "/export")
            .then(
                function(result) {
                    deferred.resolve(result.data);
                },
                function(fault) {
                    deferred.reject(fault);
                }
            );

        return deferred.promise;
    };

    this.importProfile = function(importJSON) {
        return $http.post(ENV.api['root'] + "profiles/import", importJSON)
            .then(
                function(result) {
                    messageModel.setMessages(result.data.alerts, true);
                    locationUtils.navigateToPath('/profiles/' + result.data.response.id);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

};

ProfileService.$inject = ['Restangular', '$http', '$q', 'locationUtils', 'messageModel', 'ENV'];
module.exports = ProfileService;
