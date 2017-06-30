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

var FormTenantController = function(tenant, $scope, formUtils, locationUtils, tenantService) {

    var getTenants = function() {
        tenantService.getTenants()
            .then(function(result) {
                $scope.tenants = result;
            });
    };

    $scope.tenant = tenant;

    $scope.falseTrue = [
        { value: false, label: 'false' },
        { value: true, label: 'true' }
    ];

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.hasError = formUtils.hasError;

    $scope.hasPropertyError = formUtils.hasPropertyError;

    var init = function () {
        getTenants();
    };
    init();

};

FormTenantController.$inject = ['tenant', '$scope', 'formUtils', 'locationUtils', 'tenantService'];
module.exports = FormTenantController;
