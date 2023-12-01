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

/**
 * @param {*} physLocation
 * @param {*} $scope
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../service/utils/StringUtils")} stringUtils
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/RegionService")} regionService
 */
var FormPhysLocationController = function(physLocation, $scope, $location, formUtils, stringUtils, locationUtils, regionService) {

    var getRegions = function() {
        regionService.getRegions()
            .then(function(result) {
                $scope.regions = result;
            });
    };

    $scope.physLocation = physLocation;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 100 },
        { name: 'shortName', type: 'text', required: true, maxLength: 100 },
        { name: 'address', type: 'text', required: true, maxLength: 128 },
        { name: 'city', type: 'text', required: true, maxLength: 128 },
        { name: 'state', type: 'text', required: true, maxLength: 4 },
        { name: 'zip', type: 'text', required: true, maxLength: 11 },
        { name: 'poc', type: 'text', required: false, maxLength: 128 },
        { name: 'phone', type: 'text', required: false, maxLength: 45 },
        { name: 'email', type: 'email', required: false, maxLength: 128 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.viewServers = function() {
        $location.path($location.path() + '/servers');
    };

    $scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getRegions();
    };
    init();

};

FormPhysLocationController.$inject = ['physLocation', '$scope', '$location', 'formUtils', 'stringUtils', 'locationUtils', 'regionService'];
module.exports = FormPhysLocationController;
