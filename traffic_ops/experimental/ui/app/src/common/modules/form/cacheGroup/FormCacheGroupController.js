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

var FormCacheGroupController = function(cacheGroup, $scope, $location, formUtils, locationUtils, cacheGroupService, typeService) {

    var getCacheGroups = function() {
        cacheGroupService.getCacheGroups()
            .then(function(result) {
                $scope.cacheGroups = result;
            });
    };

    var getTypes = function() {
        typeService.getTypes('cachegroup')
            .then(function(result) {
                $scope.types = result;
            });
    };

    $scope.cacheGroup = cacheGroup;

    $scope.queueUpdates = function() {
        alert('not hooked up yet: queuing updates for all cachegroup servers');
    };

    $scope.dequeueUpdates = function() {
        alert('not hooked up yet: dequeuing updates for all cachegroup servers');
    };

    $scope.viewParams = function() {
        $location.path($location.path() + '/parameters');
    };

    $scope.viewServers = function() {
        $location.path($location.path() + '/servers');
    };

    $scope.viewParams = function() {
        $location.path($location.path() + '/parameters');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getCacheGroups();
        getTypes();
    };
    init();

};

FormCacheGroupController.$inject = ['cacheGroup', '$scope', '$location', 'formUtils', 'locationUtils', 'cacheGroupService', 'typeService'];
module.exports = FormCacheGroupController;
