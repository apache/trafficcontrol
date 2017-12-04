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

var FormNewFederationController = function(cdn, federation, users, deliveryServices, $scope, $controller, federationService) {

	// extends the FormFederationController to inherit common methods
	angular.extend(this, $controller('FormFederationController', { cdn: cdn, federation: federation, deliveryServices: deliveryServices, $scope: $scope }));

	$scope.users = users;

	$scope.cname = 'New';

	$scope.settings = {
		isNew: true,
		saveLabel: 'Create'
	};

	$scope.save = function(fed) {
		federationService.createFederation(cdn, fed);
	};

};

FormNewFederationController.$inject = ['cdn', 'federation', 'users', 'deliveryServices', '$scope', '$controller', 'federationService'];
module.exports = FormNewFederationController;
