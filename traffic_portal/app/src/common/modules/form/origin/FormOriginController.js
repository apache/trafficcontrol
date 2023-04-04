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
 * This is the parent controller for all forms used to modify or create Origins.
 *
 * @param {unknown} origin
 * @param {*} $scope
 * @param {import("angular").IWindowService} $window
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../service/utils/TenantUtils")} tenantUtils
 * @param {import("../../../api/DeliveryServiceService")} deliveryServiceService
 * @param {import("../../../api/ProfileService")} profileService
 * @param {import("../../../api/TenantService")} tenantService
 * @param {import("../../../api/CoordinateService")} coordinateService
 * @param {import("../../../api/CacheGroupService")} cacheGroupService
 * @param {import("../../../models/UserModel")} userModel
 */
var FormOriginController = function(origin, $scope, $window, formUtils, locationUtils, tenantUtils, deliveryServiceService, profileService, tenantService, coordinateService, cacheGroupService, userModel) {

    var getProfiles = function() {
        profileService.getProfiles({ orderby: 'name' })
            .then(function(result) {
                $scope.profiles = result.filter(profile => profile.type === 'ORG_PROFILE');
            });
    };

    var getTenants = function() {
        tenantService.getTenant(userModel.user.tenantId)
            .then(function(tenant) {
                tenantService.getTenants()
                    .then(function(tenants) {
                        $scope.tenants = tenantUtils.hierarchySort(tenantUtils.groupTenantsByParent(tenants), tenant.parentId, []);
                        tenantUtils.addLevels($scope.tenants);
                    });
            });
    };

    var getCacheGroups = function() {
        cacheGroupService.getCacheGroups({ orderby: 'name' })
            .then(function(result) {
                $scope.cacheGroups = result.filter(cachegroup => cachegroup.typeName === 'ORG_LOC');
            });
    };

    var getCoordinates = function() {
        coordinateService.getCoordinates({ orderby: 'name' })
            .then(function(result) {
                $scope.coordinates = result;
            });
    };

    var getDeliveryServices = function() {
        deliveryServiceService.getDeliveryServices()
            .then(function(result) {
                $scope.deliveryServices = _.sortBy(result, 'xmlId');
            });
    };

    $scope.origin = origin;

    $scope.protocols = [
        { value: 'http', label: 'http' },
        { value: 'https', label: 'https' }
    ];

    $scope.nullifyIfEmptyIP = function(origin) {
        origin.ipAddress = origin.ipAddress == '' ? null : origin.ipAddress;
        origin.ip6Address = origin.ip6Address == '' ? null : origin.ip6Address;
    }

    $scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

    $scope.editDeliveryService = function(deliveryServiceId) {
        const ds = $scope.deliveryServices.find(d => d.id === deliveryServiceId);
        $window.open('/#!/delivery-services/' + ds.id + '?dsType=' + ds.type, '_blank');
    };

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getProfiles();
        getTenants();
        getCacheGroups();
        getCoordinates();
        getDeliveryServices();
    };
    init();

};

FormOriginController.$inject = ['origin', '$scope', '$window', 'formUtils', 'locationUtils', 'tenantUtils', 'deliveryServiceService', 'profileService', 'tenantService', 'coordinateService', 'cacheGroupService', 'userModel'];
module.exports = FormOriginController;
