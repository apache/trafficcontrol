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

var FormISOController = function(servers, osversions, $scope, $anchorScroll, formUtils, toolsService, messageModel) {

	$scope.servers = servers;

	$scope.osversions = osversions;

	$scope.selectedServer = {};

	$scope.falseTrue = [
		{ value: 'yes', label: 'yes' },
		{ value: 'no', label: 'no' }
	];

	$scope.iso = {
		dhcp: 'no'
	};

	$scope.isDHCP = function() {
		return $scope.iso.dhcp == 'yes';
	};

	$scope.fqdn = function(server) {
		return server.hostName + '.' + server.domainName;
	};

	$scope.copyServerAttributes = function() {
		$scope.iso = angular.extend($scope.iso, $scope.selectedServer);
	};

	$scope.generate = function(iso) {
		toolsService.generateISO(iso)
			.then(function() {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages([{level: 'success', text: 'ISO successfully downloaded'}], false);
			});
	};

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

};

FormISOController.$inject = ['servers', 'osversions', '$scope', '$anchorScroll', 'formUtils', 'toolsService', 'messageModel'];
module.exports = FormISOController;
