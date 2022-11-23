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
 * This is the controller for the form used to test consistent hashing regular
 * expressions for a DS against Traffic Router.
 *
 * @param {import("../../../api/DeliveryServiceService").DeliveryService} deliveryService
 * @param {string|RegExp} consistentHashRegex
 * @param {*} $scope
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/DeliveryServiceService")} deliveryServiceService
 */
var FormDeliveryServiceConsistentHashRegexController = function (deliveryService, consistentHashRegex, $scope, formUtils, locationUtils, deliveryServiceService) {

    $scope.deliveryService = deliveryService;

    $scope.pattern = consistentHashRegex;

    $scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    $scope.test = function (pattern, requestPath, cdnId) {
        deliveryServiceService.getConsistentHashResult(pattern, requestPath, cdnId).then(
            function (response) {
                $scope.resultingPath = response.response.resultingPathToConsistentHash;
            },
            function (response) {
                $scope.resultingPath = "ERROR GETTING RESULT FROM TRAFFIC ROUTER";
            });
    };

};

FormDeliveryServiceConsistentHashRegexController.$inject = ['deliveryService', 'consistentHashRegex', '$scope', 'formUtils', 'locationUtils', 'deliveryServiceService'];
module.exports = FormDeliveryServiceConsistentHashRegexController;
