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

/**
 * @param {*} servers
 * @param {*} osversions
 * @param {*} $scope
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../api/ToolsService")} toolsService
 * @param {import("../../../models/MessageModel")} messageModel
 * @param {import("../../../service/utils/ServerUtils")} serverUtils
 */
var FormISOController = function(servers, osversions, $scope, $anchorScroll, formUtils, toolsService, messageModel, serverUtils) {

	$scope.IPPattern = serverUtils.IPPattern;
	$scope.IPv4Pattern = serverUtils.IPv4Pattern;

	$scope.servers = servers;

	$scope.osversions = osversions;

	$scope.selectedServer = null;

	$scope.iso = {
		dhcp: false,
		interfaceMtu: 1500
	};

	$scope.fqdn = function(server) {
		return server.hostName + '.' + server.domainName;
	};

	$scope.copyServerAttributes = function() {
		const legacyNet = serverUtils.toLegacyIPInfo($scope.selectedServer.interfaces);
		$scope.iso.hostName = $scope.selectedServer.hostName;
		$scope.iso.domainName = $scope.selectedServer.domainName;
		$scope.iso.interfaceName = legacyNet.interfaceName;
		$scope.iso.interfaceMtu = legacyNet.interfaceMtu;
		$scope.iso.ip6Address = legacyNet.ip6Address;
		$scope.iso.ip6Gateway = legacyNet.ip6Gateway;
		$scope.iso.ipAddress = legacyNet.ipAddress;
		$scope.iso.ipGateway = legacyNet.ipGateway;
		$scope.iso.ipNetmask = legacyNet.ipNetmask;
		$scope.iso.mgmtIpAddress = $scope.selectedServer.mgmtIpAddress;
		$scope.iso.mgmtIpNetmask = $scope.selectedServer.mgmtIpNetmask;
		$scope.iso.mgmtIpGateway = $scope.selectedServer.mgmtIpGateway;
		$scope.iso.mgmtInterface = $scope.selectedServer.mgmtInterface;
	};

	$scope.generate = function(iso) {
		// for whatever reason this was designed with "yes" and "no" instead of actual
		// boolean values, so we need to emulate that here.
		const tmp = Object.assign({}, iso);
		tmp.dhcp = iso.dhcp ? "yes" : "no";
		toolsService.generateISO(tmp)
			.then(function() {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages([{level: 'success', text: 'ISO successfully downloaded'}], false);
			});
	};

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

};

FormISOController.$inject = ['servers', 'osversions', '$scope', '$anchorScroll', 'formUtils', 'toolsService', 'messageModel', "serverUtils"];
module.exports = FormISOController;
