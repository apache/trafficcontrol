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

var FormISOController = function(servers, osversions, $scope, $anchorScroll, formUtils, toolsService, messageModel, serverUtils) {

	// This RegExp matches any valid IPv4 or IPv6 address
	// Source: Rahul Tripathy @ https://stackoverflow.com/questions/32324614/how-to-validate-ipv6-address-in-angularjs#answer-32324868
	// May Neptune have mercy on my soul...
	$scope.IPPattern = /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\:){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?|((25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)$/

	// ... this one just matches valid IPv4 addresses (or netmasks)
	$scope.IPv4Pattern = /^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)$/

	$scope.servers = servers;

	$scope.osversions = osversions;

	$scope.selectedServer = null;

	$scope.falseTrue = [
		{ value: 'yes', label: 'yes' },
		{ value: 'no', label: 'no' }
	];

	$scope.iso = {
		dhcp: false,
		interfaceMtu: 1500
	};

	$scope.isDHCP = function() {
		return $scope.iso.dhcp == 'yes';
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
		iso.dhcp = iso.dhcp ? "yes" : "no";
		toolsService.generateISO(iso)
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
