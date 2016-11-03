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

var TableCacheGroupsController = function(cacheGroups, $scope, locationUtils) {

    $scope.cacheGroups = cacheGroups;

    $scope.editCacheGroup = function(id) {
        locationUtils.navigateToPath('/configure/cache-groups/' + id + '/edit');
    };

    $scope.createCacheGroup = function() {
        locationUtils.navigateToPath('/configure/cache-groups/new');
    };

    angular.element(document).ready(function () {
        $('#cacheGroupsTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": -1
        });
    });

};

TableCacheGroupsController.$inject = ['cacheGroups', '$scope', 'locationUtils'];
module.exports = TableCacheGroupsController;
