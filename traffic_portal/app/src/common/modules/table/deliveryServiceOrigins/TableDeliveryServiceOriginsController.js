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
 * @param {*} deliveryService
 * @param {*} origins
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 */
var TableDeliveryServiceOriginsController = function(deliveryService, origins, $scope, $state, locationUtils) {

    $scope.createOrigin = function() {
        var path = '/origins/new';
        locationUtils.navigateToPath(path);
    };

    $scope.deliveryService = deliveryService;

    $scope.origins = origins;

    $scope.editOrigin = function(id) {
        var path = '/origins/' + id;
        locationUtils.navigateToPath(path);
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    $scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

    angular.element(document).ready(function () {
        $('#originsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "aaSorting": []
        });
    });

};

TableDeliveryServiceOriginsController.$inject = ['deliveryService', 'origins', '$scope', '$state', 'locationUtils'];
module.exports = TableDeliveryServiceOriginsController;
