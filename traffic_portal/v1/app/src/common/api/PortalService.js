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

var PortalService = function($http, $q) {

    this.getReleaseVersionInfo = function() {
        var deferred = $q.defer();
        $http.get('traffic_portal_release.json')
            .success(function(result) {
                deferred.resolve(result);
            });

        return deferred.promise;
    };

    this.getProperties = function() {
        var deferred = $q.defer();
        $http.get('traffic_portal_properties.json')
            .success(function(result) {
                deferred.resolve(result.properties);
            });

        return deferred.promise;
    };

};

PortalService.$inject = ['$http', '$q'];
module.exports = PortalService;