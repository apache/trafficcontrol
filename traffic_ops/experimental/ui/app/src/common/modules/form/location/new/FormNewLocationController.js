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

var FormNewLocationController = function(location, $scope, $controller, locationService) {

    // extends the FormLocationController to inherit common methods
    angular.extend(this, $controller('FormLocationController', { location: location, $scope: $scope }));

    $scope.locationName = 'New';

    $scope.settings = {
        isNew: true,
        saveLabel: 'Create'
    };

    $scope.save = function(location) {
        locationService.createLocation(location);
    };

};

FormNewLocationController.$inject = ['location', '$scope', '$controller', 'locationService'];
module.exports = FormNewLocationController;
