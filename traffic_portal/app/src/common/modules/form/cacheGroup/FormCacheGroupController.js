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
        cacheGroupService.getCacheGroups({ orderby: 'name' })
            .then(function(result) {
                $scope.cacheGroups = result;
            });
    };

    var getTypes = function() {
        typeService.getTypes({ useInTable: 'cachegroup' })
            .then(function(result) {
                $scope.types = result;
            });
    };

    $scope.cacheGroup = cacheGroup;

    $scope.viewAsns = function() {
        $location.path($location.path() + '/asns');
    };

    $scope.viewParams = function() {
        $location.path($location.path() + '/parameters');
    };

    $scope.viewServers = function() {
        $location.path($location.path() + '/servers');
    };

    $scope.viewStaticDnsEntries = function() {
        $location.path($location.path() + '/static-dns-entries');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    $scope.localizationMethods = {
        DEEP_CZ: false,
        CZ: false,
        GEO: false
    };

    $scope.setLocalizationMethods = function(cacheGroup) {
        var methods = [];
        var keys = Object.keys($scope.localizationMethods);
        for (var i = 0; i < keys.length; i++) {
            if ($scope.localizationMethods[keys[i]]) {
                methods.push(keys[i]);
            }
        }
        cacheGroup.localizationMethods = methods;
    };

    var initLocalizationMethods = function() {
        // by default, no explicitly enabled methods means ALL are enabled
        if (!cacheGroup.localizationMethods) {
            var keys = Object.keys($scope.localizationMethods);
            for (var i = 0; i < keys.length; i++) {
                $scope.localizationMethods[keys[i]] = true;
            }
            return;
        }
        for (var i = 0; i < cacheGroup.localizationMethods.length; i++) {
            if ($scope.localizationMethods.hasOwnProperty(cacheGroup.localizationMethods[i])) {
                $scope.localizationMethods[cacheGroup.localizationMethods[i]] = true;
            }
        }
    };

    var init = function () {
        initLocalizationMethods();
        getCacheGroups();
        getTypes();
    };
    init();

};

FormCacheGroupController.$inject = ['cacheGroup', '$scope', '$location', 'formUtils', 'locationUtils', 'cacheGroupService', 'typeService'];
module.exports = FormCacheGroupController;
