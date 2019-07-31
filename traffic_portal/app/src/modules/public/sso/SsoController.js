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

var SsoController = function($scope, $location, authService, propertiesModel) {

	var init = function () {
        const authCodeTokenUrl = propertiesModel.properties.oAuth.oAuthCodeTokenUrl;
        const clientId = propertiesModel.properties.oAuth.clientId;
        const redirectUri = localStorage.getItem('redirectUri');

        var code = '';
        if ($location.hash()) {
            code = $location.hash().replace('code=', '');
        } else {
            code = $location.search()['code'];
        }
        authService.oauthLogin(authCodeTokenUrl, code, clientId, redirectUri);
	};
	init();

};

SsoController.$inject = ['$scope', '$location', 'authService', 'propertiesModel'];
module.exports = SsoController;
