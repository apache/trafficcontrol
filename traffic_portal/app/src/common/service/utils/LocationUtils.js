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

var LocationUtils = function($location, $uibModal) {

    this.navigateToPath = function(path, unsavedChanges) {
        if (unsavedChanges) {
            const params = {
                title: 'Confirm Navigation',
                message: 'You have unsaved changes that will be lost if you decide to continue.<br><br>Do you want to continue?'
            };
            let modalInstance = $uibModal.open({
                templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
                controller: 'DialogConfirmController',
                size: 'md',
                resolve: {
                    params: function () {
                        return params;
                    }
                }
            });
            modalInstance.result.then(function() {
                $location.url(path);
            });
        } else {
            $location.url(path);
        }
    };

};

LocationUtils.$inject = ['$location', '$uibModal'];
module.exports = LocationUtils;
