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
 * @typedef LegacyIPInfo
 * @property {string | null} ipAddress
 * @property {string | null} ipGateway
 * @property {string | null} ipNetmask
 * @property {string | null} ip6Address
 * @property {string | null} ip6Gateway
 * @property {string | null} interfaceName
 * @property {number | null} interfaceMtu
 * @property {string | null} routerHostName
 * @property {string | null} routerPortName
 */

/**
 * @typedef ServerIP
 * @property {string} address
 * @property {boolean} serviceAddress
 * @property {string | null} gateway
 */

/**
 * @typedef ServerInterface
 * @property {string} name
 * @property {ServerIP[]} ipAddresses
 * @property {number | null} mtu
 * @property {string | null} routerHostName
 * @property {string | null} routerPortName
 */

/**
 * ServerUtils provides methods for dealing with servers of different types (and
 * Types).
 */
class ServerUtils {

	// This RegExp matches any valid IPv4 or IPv6 address - CIDRs allowed on IPv6 addresses only
	// Source: Rahul Tripathy @ https://stackoverflow.com/questions/32324614/how-to-validate-ipv6-address-in-angularjs#answer-32324868
	// May Neptune have mercy on my soul...
	IPPattern = /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\:){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?(\/(12[0-8]|1[0-1][0-9]|[0-9][0-9]?))?$|^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)$/;
	// ... this one allows IPv4 addresses to have CIDR-notation subnets
	IPWithCIDRPattern = /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\:){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?(\/(12[0-8]|1[0-1][0-9]|[0-9][0-9]?))?$|^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)(\/(3[0-2]|[0-2]?[0-9]))?$/;
	// ... this one just matches valid IPv4 addresses (or netmasks)
	IPv4Pattern = /^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)$/;

	/**
	 * @param {import("angular").IWindowService} $window
	 * @param {import("../../models/PropertiesModel")} propertiesModel
	 * @param {import("../../models/UserModel")} userModel
	 */
	constructor($window, propertiesModel, userModel) {
		/** @private */
		this.$window = $window;
		/** @private */
		this.propertiesModel = propertiesModel;
		/** @private */
		this.userModel = userModel;
	}

	/**
	 * Checks if a server is a cache server.
	 *
	 * @param {{type?: string | null | undefined}} server
	 * @returns {boolean}
	 */
	isCache(server) {
		return !!server.type && (server.type.startsWith("EDGE") || server.type.startsWith("MID"));
	}

	/**
	 * Checks if a server is an edge-tier cache server.
	 *
	 * @param {{type?: string | null | undefined}} server
	 * @returns {boolean}
	 */
	isEdge(server) {
		return !!server.type && server.type.startsWith("EDGE");
	}

	/**
	 * Checks if a server is an origin server.
	 *
	 * @param {{type?: string | null | undefined}} server
	 * @returns {boolean}
	 */
	isOrigin(server) {
		return !!server.type && server.type.startsWith("ORG");
	}

	/**
	 * Checks if a Status is considered by health protocol to be "offline". Note
	 * that any unrecognized Status is also thusly treated and that isn't
	 * checked by this method.
	 *
	 * @param {string} status
	 * @returns {boolean}
	 */
	isOffline(status) {
		return (status === "OFFLINE" || status === "ADMIN_DOWN");
	}

	/**
	 * Gets a server's offline reason - or "None" if it doesn't have one.
	 *
	 * @param {{offlineReason?: string | null | undefined}} server
	 * @returns {string}
	 */
	offlineReason(server) {
		return (server.offlineReason) ? server.offlineReason : "None";
	}

	/**
	 * Redirects the current browsing context to an SSH URL using the currently
	 * authenticated user's username and the provided IP address.
	 *
	 * @deprecated Links should be links; use ng-href with the appropriate link
	 * on an Anchor element instead of utilizing this method in a click handler.
	 *
	 * @param {string} ip
	 * @param {Event} $event
	 */
	ssh(ip, $event) {
		if (ip && ip.length > 0) {
			this.$window.location.href = `ssh://${this.userModel.user.username}@${ip}`;
		}
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
	}

	/**
	 * Opens a custom charts URL for the given server in a new browsing context.
	 *
	 * Note that this will break in half if a custom charts base URL is not
	 * configured for servers, and the caller is expected to assume that
	 * responsibility.
	 *
	 * @deprecated Links should be links; use ng-href with the appropriate link
	 * on an Anchor element instead of utilizing this method in a click handler.
	 *
	 * @param {{hostName: string}} server
	 * @param {Event} $event
	 */
	openCharts(server, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}
		this.$window.open(
			this.propertiesModel.properties.servers.charts.baseUrl + server.hostName,
			"_blank"
		);
	}

	/**
	 * Converts a server's interfaces into legacy IP information. (primarily
	 * for use in tables)
	 *
	 * It does this by returning only the service addresses of the server.
	 *
	 * @param {Array<ServerInterface>} interfaces - The interfaces of the server
	 * to be converted.
	 * @returns {LegacyIPInfo} An object with all of the legacy properties of
	 * non-interface-based servers.
	 */
	toLegacyIPInfo(interfaces) {
		/** @type LegacyIPInfo */
		const legacyInfo = {
			ipAddress: null,
			ipGateway: null,
			ipNetmask: null,
			ip6Address: null,
			ip6Gateway: null,
			interfaceName: null,
			interfaceMtu: null,
			routerHostName: null,
			routerPortName: null
		};
		if (!interfaces) {
			return legacyInfo;
		}

		for (let i = 0; i < interfaces.length; ++i) {
			const inf = interfaces[i];

			for (let j = 0; j < inf.ipAddresses.length; ++j) {
				const ip = inf.ipAddresses[j];
				if (!ip.serviceAddress) {
					continue;
				}
				legacyInfo.interfaceName = inf.name;
				legacyInfo.interfaceMtu = inf.mtu;
				legacyInfo.routerHostName = inf.routerHostName;
				legacyInfo.routerPortName = inf.routerPortName;
				let address = ip.address;

				// we don't validate ips here; if it has a '.' it's ipv4,
				// otherwise it's ipv6
				if (address.includes(".")) {
					if (address.includes("/")) {
						const parts = address.split("/");
						address = parts[0];
						let masklen = Number(parts[1]);

						const mask = [];
						for (let k = 0; k < 4; ++k) {
							const n = Math.min(masklen, 8);
							mask.push(256 - Math.pow(2, 8 - n));
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
	};

}

ServerUtils.$inject = ["$window", "propertiesModel", "userModel"];
module.exports = ServerUtils;
