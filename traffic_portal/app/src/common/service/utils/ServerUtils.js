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

var ServerUtils = function($window, propertiesModel, userModel) {

	this.isCache = function(server) {
		return server.type && (server.type.indexOf('EDGE') == 0 || server.type.indexOf('MID') == 0);
	};

	this.isEdge = function(server) {
		return server.type && (server.type.indexOf('EDGE') == 0);
	};

	this.isOrigin = function(server) {
		return server.type && (server.type.indexOf('ORG') == 0);
	};

	this.isOffline = function(status) {
		return (status == 'OFFLINE' || status == 'ADMIN_DOWN');
	};

	this.offlineReason = function(server) {
		return (server.offlineReason) ? server.offlineReason : 'None';
	};

	this.ssh = function(ip, $event) {
		if (ip && ip.length > 0) {
			$window.location.href = 'ssh://' + userModel.user.username + '@' + ip;
		}
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
	};

	this.openCharts = function(server, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}
		$window.open(
			propertiesModel.properties.servers.charts.baseUrl + server.hostName,
			'_blank'
		);
	};

	/**
	 * Converts a server's interfaces into legacy IP information. (primarily
	 * for use in tables)
	 *
	 * It does this by returning only the service addresses of the server.
	 *
	 * @param {Array<object>} interfaces - The interfaces of the server to be converted
	 * @returns {object} An object with all of the legacy properties of non-
	 * interface-based servers: ipAddress, ipGateway, ipNetmask, ip6Address,
	 * ip6Gateway, interfaceName, and interfaceMtu
	 */
	this.toLegacyIPInfo = function(interfaces) {
		const legacyInfo = {
			ipAddress: null,
			ipGateway: null,
			ipNetmask: null,
			ip6Address: null,
			ip6Gateway: null,
			interfaceName: null,
			interfaceMtu: null
		};
		for (const inf of interfaces) {
			legacyInfo.interfaceName = inf.name;
			legacyInfo.interfaceMtu = inf.mtu;

			for (const ip of inf.ipAddresses) {
				if (!ip.serviceAddress) {
					continue;
				}

				const address = ip.address;

				// we don't validate ips here; if it has a '.' it's ipv4,
				// otherwise it's ipv6
				if (address.contains(".")) {
					if (address.contains("/")) {
						const parts = address.split("/");
						address = parts[0];
						let masklen = Number(parts[1]);

						const mask = [];
						for (let i = 0; i < 4; ++i) {
							const n = Math.min(masklen, 8);
							mask.push(256 - Math.pow(2, 8-n));
							masklen -= n;
						}
						legacyInfo.ipNetmask = mask.join(".");
					}
					legacyInfo.ipAddress = address;
					legacyInfo.ipGateway = ip.gateway;
				} else {
					legacyInfo.ip6Address = address;
					legacyInfo.ip6Gateway = ip.gateway;
				}

				if (legacyInfo.ipAddress && legacyInfo.ip6Address) {
					break;
				}
			}

			if (legacyInfo.ipAddress && legacyInfo.ip6Address) {
				break;
			}
		}

		return legacyInfo;
	}

};

ServerUtils.$inject = ['$window', 'propertiesModel', 'userModel'];
module.exports = ServerUtils;
