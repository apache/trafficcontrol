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

var UserModel = function($rootScope, propertiesModel) {

    this.loaded = false;

    this.user = {};

    this.setUser = function(userData) {
        this.loaded = true;
        this.user = userData;
        $rootScope.$broadcast('userModel::userUpdated', this.user);
    };

    this.resetUser = function() {
        this.loaded = false;
        this.userId = 0;
        this.user = {};
        $rootScope.$broadcast('userModel::userUpdated', this.user);
    };

    this.hasCapability = function(cap) {
        if (propertiesModel.properties.enforceCapabilities == false) {
            return true;
        }
        return _.has(this.user, 'capabilities') && _.indexOf(this.user.capabilities, cap) != -1;
    };

};

UserModel.$inject = ['$rootScope', 'propertiesModel'];
module.exports = UserModel;
