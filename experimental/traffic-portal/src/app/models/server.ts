/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

/**
 * This file is for modeling and functionality related to Server objects
 */

/** IPAddress is a single IP address of a single network interface of a server. */
export interface IPAddress {
	/** The actual IP address. */
	address: string;
	/** The IP address of a gateway for this IP, if one exists/is known. */
	gateway: string | null;
	/** Whether or not this IP address is responsible for serving ATC traffic. */
	serviceAddress: boolean;
}

/** Interface is a server's network interface. */
export interface Interface {
	/** The IP addresses assigned to this network interface. */
	ipAddresses: Array<IPAddress>;
	/** The maximum bandwidth for considering the server healthy, if any. */
	maxBandwidth: number | null;
	/**
	 * Whether or not the Traffic Monitor should consider this network interface
	 * in health calculations.
	 */
	monitor: boolean;
	/** The maximum transmission unit of the network interface, if known. */
	mtu: number | null;
	/** The name of the network interface in the `/dev` directory. */
	name: string;
}

/**
 * Server represents a nebulous "server" object in Traffic Ops.
 */
export interface Server {
	/** The Cache Group to which the server belongs. */
	cachegroup?: string;
	/** The integral, unique identifier of the Cache Group to which the server belongs. */
	cachegroupId: number;
	/** The integral, unique identifier of the CDN to which the server belongs. */
	cdnId: number;
	/** The name of the CDN to which the server belongs. */
	cdnName?: string;
	/**
	 * The servers FQDN without it's hostname - e.g. 'apache.org' from
	 * 'trafficcontrol.apache.org'.
	 */
	domainName: string;
	/** legacy field with no purpose. */
	guid?: number;
	/**
	 * The server's hostname, e.g. 'trafficcontrol' from
	 * 'trafficcontrol.apache.org'.
	 */
	hostName: string;
	/** The port used to serve HTTPS responses, if any. */
	httpsPort: number | null;
	/** An integral, unique identifier for this Server. */
	id?: number;
	/** The IP address of the Server's ILO interface. */
	iloIpAddress: string | null;
	/** The IP address of the gateway to the Server's ILO interface. */
	iloIpGateway: string | null;
	/** A netmask that describes the subnet allocated to the Server's ILO interface. */
	iloIpNetmask: string | null;
	/** The Server's ILO interface's password. */
	iloPassword: string | null;
	/** The Server's ILO interface's root user's name. */
	iloUsername: string | null;
	/** The Server's network interfaces. */
	interfaces: Array<Interface>;
	/** The date/time at which the Server was last updated. */
	lastUpdated?: Date;
	/** The IP address of the server's management interface. */
	mgmtIpAddress: string | null;
	/** The IP address of the gateway to the Server's management interface. */
	mgmtIpGateway: string | null;
	/** The netmask that describes the subnet allocated to the Server's management interface. */
	mgmtIpNetmask: string | null;
	/** The reason the Server has been taken out of service. */
	offlineReason: string | null;
	/** The physical location in which the Server resides. */
	physLocation?: string;
	/** An integral, unique identifier for the physical location in which the Server resides. */
	physLocationId: number;
	/** The Profile used by the Server. */
	profile?: string;
	/** A description of the Profile used by the Server. */
	profileDesc?: string | null;
	/** An integral, unique identifier for the Profile used by the Server. */
	profileId: number;
	/** Whether or not revalidations are pending for this Server. */
	revalPending: boolean;
	/** legacy field with no purpose. */
	rack?: string | null;
	/** The hostname of the router that routes to this Server. */
	routerHostName?: string | null;
	/** The... name... of the port... used by the Server's router?? */
	routerPortName?: string | null;
	/** The Server's status. */
	status?: string;
	/** An integral, unique, identifier for the Server's Status. */
	statusId: number;
	/** The time at which the server's status was last updated. */
	statusLastUpdated?: Date | null;
	/** The port on which the Server listens for incoming TCP connections. */
	tcpPort: number | null;
	/** The type of the Server. */
	type?: string;
	/** An integral, unique identifier for the Type of this Server. */
	typeId: number;
	/** Whether or not this Server has updates pending. */
	updPending: boolean;
	/** The string used by Traffic Router for consistent hashing to this Server. */
	xmppId?: string | null;
	/** legacy field with no purpose. */
	xmppPasswd?: string | null;
}

export const defaultServer: Server = {
	cachegroupId: -1,
	cdnId: -1,
	domainName: "",
	hostName: "",
	httpsPort: null,
	iloIpAddress: null,
	iloIpGateway: null,
	iloIpNetmask: null,
	iloPassword: null,
	iloUsername: null,
	interfaces: [],
	mgmtIpAddress: null,
	mgmtIpGateway: null,
	mgmtIpNetmask: null,
	offlineReason: null,
	physLocationId: -1,
	profileId: -1,
	revalPending: false,
	statusId: -1,
	tcpPort: null,
	typeId: -1,
	updPending: false
};

/**
 * Servercheck models the data returned by the /servercheck API endpoint.
 */
export interface Servercheck {
	/** contains the server's Status */
	adminState: string;
	/** the name of the Cache Group to which the server belongs */
	cacheGroup: string;
	/**
	 * Checks emulates a map of check names to their numbers. All values are
	 * numbers, but some may express boolean concepts; for example, the ORT
	 * check uses 1 to represent "true" and any other value indicates "false".
	 */
	checks?: Record<string, number>;
	/** the server's hostname */
	hostName: string;
	/** the server's ID */
	id: number;
	/** the name of the server's Profile */
	profile: string;
	/** whether or not the server has pending revalidations */
	revalPending: boolean;
	/** the name of the server's Type */
	type: string;
	/** whether or not the server has updates pending */
	updPending: boolean;
}

/**
 * Builds a true Map from the Servercheck's "checks" property.
 *
 * @param srv The Servercheck to convert.
 * @returns A map of servercheck check names to their values.
 */
export function checkMap(srv: Servercheck): Map<string, number | boolean> {
	const ret = new Map();
	if (!srv.checks) {
		return ret;
	}
	for (const [key, value] of Object.entries(srv.checks)) {
		switch (key) {
			case "ILO":
			case "10G":
			case "FQDN":
			case "DSCP":
			case "10G6":
			case "MTU":
				ret.set(key, value === 1);
				break;
			default:
				ret.set(key, value);
				break;
		}
	}
	return ret;
}

export const DUMMY_SERVER: Server = {
	cachegroupId: -1,
	cdnId: -1,
	domainName: "",
	hostName: "",
	httpsPort: null,
	iloIpAddress: null,
	iloIpGateway: null,
	iloIpNetmask: null,
	iloPassword: null,
	iloUsername: null,
	interfaces: [],
	mgmtIpAddress: null,
	mgmtIpGateway: null,
	mgmtIpNetmask: null,
	offlineReason: null,
	physLocationId: -1,
	profileId: -1,
	revalPending: false,
	statusId: -1,
	tcpPort: null,
	typeId: -1,
	updPending: false
};
