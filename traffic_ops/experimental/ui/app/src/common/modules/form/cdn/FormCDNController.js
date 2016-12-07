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

var FormCDNController = function(cdn, $scope, formUtils, stringUtils, locationUtils) {

    $scope.cdn = cdn;

    $scope.props = [
        { name: 'name', type: 'text', required: true, maxLength: 45 }
    ];

    $scope.labelize = stringUtils.labelize;

    $scope.falseTrue = [
        { value: false, label: 'false' },
        { value: true, label: 'true' }
    ];

    $scope.manageDNSSEC = function() {
        alert('not hooked up yet: manageDNSSEC for CDN');
    };

    $scope.manageSSL = function() {
        alert('not hooked up yet: manageSSL for cdn');
    };

    $scope.cachegroupHealth = function() {
        alert('not hooked up yet: cachegroupHealth for CDN');
    };

    $scope.queueUpdates = function() {
        alert('not hooked up yet: queuing updates for all cdn servers');
    };

    $scope.dequeueUpdates = function() {
        alert('not hooked up yet: dequeuing updates for all cdn servers');
    };

    $scope.manageSnapshots = function() {
        alert('not hooked up yet: manageSnapshots for CDN');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

};

FormCDNController.$inject = ['cdn', '$scope', 'formUtils', 'stringUtils', 'locationUtils'];
module.exports = FormCDNController;
