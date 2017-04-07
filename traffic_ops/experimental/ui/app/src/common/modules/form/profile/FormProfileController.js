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

var FormProfileController = function(profile, $scope, $location, formUtils, locationUtils, cdnService) {

    var getCDNs = function() {
        cdnService.getCDNs()
            .then(function(result) {
                $scope.cdns = result;
            });
    };

    $scope.profile = profile;

    $scope.types = [
        { value: 'ATS_PROFILE' },
        { value: 'TR_PROFILE' },
        { value: 'TM_PROFILE' },
        { value: 'TS_PROFILE' },
        { value: 'TP_PROFILE' },
        { value: 'INFLUXDB_PROFILE' },
        { value: 'RIAK_PROFILE' },
        { value: 'SPLUNK_PROFILE' },
        { value: 'DS_PROFILE' },
        { value: 'ORG_PROFILE' },
        { value: 'KAFKA_PROFILE' },
        { value: 'LOGSTASH_PROFILE' },
        { value: 'ES_PROFILE' },
        { value: 'UNK_PROFILE' }
    ];

    $scope.viewParams = function() {
        $location.path($location.path() + '/parameters');
    };

    $scope.viewServers = function() {
        $location.path($location.path() + '/servers');
    };

    $scope.viewDeliveryServices = function() {
        $location.path($location.path() + '/delivery-services');
    };

    $scope.cloneProfile = function() {
        alert('not hooked up yet: cloneProfile');
    };

    $scope.exportProfile = function() {
        alert('not hooked up yet: exportProfile');
    };

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getCDNs();
    };
    init();

};

FormProfileController.$inject = ['profile', '$scope', '$location', 'formUtils', 'locationUtils', 'cdnService'];
module.exports = FormProfileController;
