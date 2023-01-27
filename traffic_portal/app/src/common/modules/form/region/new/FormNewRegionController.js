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
 * @param {*} region
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/RegionService")} regionService
 */
var FormNewRegionController = function(region, $scope, $controller, locationUtils, regionService) {

    // extends the FormRegionController to inherit common methods
    angular.extend(this, $controller('FormRegionController', { region: region, $scope: $scope }));

    $scope.regionName = 'New';

    $scope.settings = {
        isNew: true,
        saveLabel: 'Create'
    };

    $scope.save = function(region) {
        regionService.createRegion(region).
            then(function() {
                locationUtils.navigateToPath('/regions');
            });
    };

};

FormNewRegionController.$inject = ['region', '$scope', '$controller', 'locationUtils', 'regionService'];
module.exports = FormNewRegionController;
