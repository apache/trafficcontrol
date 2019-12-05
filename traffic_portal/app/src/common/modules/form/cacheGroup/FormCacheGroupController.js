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

var FormCacheGroupController = function(cacheGroup, types, cacheGroups, $scope, $location, formUtils, locationUtils, cacheGroupService) {

    $scope.types = types;

    $scope.cacheGroups = cacheGroups;

    $scope.cacheGroup = cacheGroup;

    $scope.cacheGroupFallbackUpdated = false;

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

    $scope.cacheGroupFallbackOptions = [];

    $scope.selectedCacheGroupFallbackOptions = [];

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

    var initCacheGroupFallbackGeo = function() {
        if (cacheGroup.fallbackToClosest == null || cacheGroup.fallbackToClosest === '') {
            cacheGroup.fallbackToClosest = true;
        }
    };

    function CacheGroupFallbackOption(index, group) {
        this.index = index;
        this.group = group;
    }

    // Creates a list of available Fallback options and a list of previously selected Fallback options
    $scope.getFallbackOptions = function() {
        for (var i = 0; i < $scope.cacheGroups.length; i++) {
            var cg = $scope.cacheGroups[i];
            // Fallbacks are required to be of type EDGE_LOC and a cachegroup cannot fallback to itself so these are skipped and the loop is continued
            if (!$scope.isEdgeLoc(cg.typeId) || cg.name == cacheGroup.name) continue;
            var fb = new CacheGroupFallbackOption(i, $scope.cacheGroups[i].name);
            // If the fallback has not been used yet, it is added to the list of available fallbacks
            if (cacheGroup.fallbacks == null || cacheGroup.fallbacks.indexOf(cg.name) < 0) {
                $scope.cacheGroupFallbackOptions.push(fb);
            } else {
                // If fallback has been selected previously then it is added to the list of selected fallbacks
                $scope.selectedCacheGroupFallbackOptions.push(fb);
            }
        }
    };

    $scope.fallbackSelected = '';

    $scope.draggedFallback = '';

    $scope.droppedFallback = '';

    $scope.moveAbove = true;

    // Updates the list of already selected fallbacks and removes the newly selected fallback from the list of available fallbacks
    $scope.updateFallbacks = function(cacheGroup) {
        if (cacheGroup.fallbacks == null) {
            cacheGroup.fallbacks = new Array();
        }
        // Add selected fallback to selected list if it is not already there
        if ($scope.fallbackSelected && cacheGroup.fallbacks.indexOf($scope.fallbackSelected) === -1) {
            cacheGroup.fallbacks.push($scope.fallbackSelected);
        }
        // Update list of available fallbacks so it does not include the newly selected fallback
        for (var i = 0; i < $scope.cacheGroupFallbackOptions.length; i++) {
            var fbo = $scope.cacheGroupFallbackOptions[i];
            if (fbo.group === $scope.fallbackSelected) {
                // Removes selected fallback from list of availables
                $scope.cacheGroupFallbackOptions.splice($scope.cacheGroupFallbackOptions.indexOf(fbo), 1);
                // Adds selected fallback to list of selected
                $scope.selectedCacheGroupFallbackOptions.push(fbo);
                break;
            }
        }
        $scope.fallbackSelected = '';
    };

    $scope.updateForNewType = function() {
        // removes Cache Group fallbacks if type has changed and is no longer EDGE_LOC
        if (!$scope.isEdgeLoc(cacheGroup.typeId)) {
            let currentFallbacksCount = cacheGroup.fallbacks.length;
            for (var i = 0; i < currentFallbacksCount; i++) {
                // removes fallbacks at position 0 since array is changing every loop
                $scope.removeFallback(cacheGroup.fallbacks[0]);
            }
        }
    };

    $scope.save = function(cacheGroup) {
        $scope.setLocalizationMethods(cacheGroup);
        cacheGroupService.createCacheGroup(cacheGroup);
        $scope.cacheGroupFallbackUpdated = false;
    };

    $scope.removeFallback = function(fb) {
        cacheGroup.fallbacks.splice(cacheGroup.fallbacks.indexOf(fb), 1);
        $scope.cacheGroupFallbackUpdated = true;
        for (var i = 0; i < $scope.selectedCacheGroupFallbackOptions.length; i++) {
            var fbo = $scope.selectedCacheGroupFallbackOptions[i];
            if (fbo.group === fb) {
                $scope.selectedCacheGroupFallbackOptions.splice($scope.selectedCacheGroupFallbackOptions.indexOf(fbo), 1);
                for (var j = 0; j < $scope.cacheGroupFallbackOptions.length; j++) {
                    if ($scope.cacheGroupFallbackOptions[j].index > fbo.index) {
                        $scope.cacheGroupFallbackOptions.splice(j, 0, fbo);
                        break;
                    } else if (j === $scope.cacheGroupFallbackOptions.length - 1) {
                        $scope.cacheGroupFallbackOptions.splice(j + 1, 0, fbo);
                        break;
                    }
                }
                break;
            }
        }
    };

    $scope.handleDrag = function(fb) {
        $scope.draggedFallback = fb;
    };

    $scope.handleDrop = function(fb) {
        $scope.droppedFallback = fb;
        var draggedIndex = cacheGroup.fallbacks.indexOf($scope.draggedFallback);
        var droppedIndex = cacheGroup.fallbacks.indexOf($scope.droppedFallback);
        var newIndex = droppedIndex;
        if (draggedIndex < droppedIndex) {
            newIndex = droppedIndex - 1;
        }
        if (!$scope.moveAbove) {
            newIndex = newIndex + 1;
        }
        cacheGroup.fallbacks.splice(draggedIndex, 1);
        cacheGroup.fallbacks.splice(newIndex, 0, $scope.draggedFallback);
        $scope.cacheGroupFallbackUpdated = true;
    };

    $scope.isEdgeLoc = function(id) {
        var selectedType = '';
        if ($scope.types != null) {
            for (var i = 0; i < $scope.types.length; i++) {
                if ($scope.types[i].id == id) {
                    selectedType = $scope.types[i].name;
                    break;
                }
            }
        }
        return selectedType == 'EDGE_LOC';
    };

    var init = function () {
        initLocalizationMethods();
        $scope.getFallbackOptions();
        initCacheGroupFallbackGeo();
    };
    init();
};

FormCacheGroupController.$inject = ['cacheGroup', 'types', 'cacheGroups', '$scope', '$location', 'formUtils', 'locationUtils', 'cacheGroupService', 'typeService'];
module.exports = FormCacheGroupController;
