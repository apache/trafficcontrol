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

export interface IPAddress {
	address: string;
	gateway: string | null;
	serviceAddress: boolean;
}

export interface Interface {
	ipAddresses: Array<IPAddress>;
	maxBandwidth: number | null;
	monitor: boolean;
	mtu: number | null;
	name: string;
}

export interface Server {
	cachegroup?: string;
	cachegroupId: number;
	cdnId: number;
	cdnName?: string;
	domainName: string;
	guid?: number;
	hostName: string;
	httpsPort: number | null;
	id: number | null;
	iloIpAddress: string | null;
	iloIpGateway: string | null;
	iloIpNetmask: string | null;
	iloPassword: string | null;
	iloUsername: string | null;
	interfaces: Array<Interface>;
	lastUpdated?: Date;
	mgmtIpAddress: string | null;
	mgmtIpGateway: string | null;
	mgmtIpNetmask: string | null;
	offlineReason: string | null;
	physLocation?: string;
	physLocationId: number;
	profile?: string;
	profileDesc?: string | null;
	profileId: number;
	revalPending: boolean;
	rack?: string | null;
	routerHostName?: string | null;
	routerPortName?: string | null;
	status?: string;
	statusId: number;
	tcpPort: number | null;
	type?: string;
	typeId: number;
	updPending: boolean;
	xmppId?: string | null;
	xmppPasswd?: string | null;

}

export interface Servercheck {
	profile: string;
	cacheGroup: string;
	checks?: {[key:string]: number};
	updPending: boolean;
	revalPending: boolean;
	hostName: string;
	adminState: string;
	id: number;
	type: string;

	checkMap(): Map<string, number|boolean>;
}

export function checkMap(this: Servercheck): Map<string, number|boolean> {
	const ret = new Map();
	if (!this.checks) {
		return ret;
	}
	for (const [key, value] of Object.entries(this.checks)) {
		switch (key) {
			case "ILO":
			case "10G":
			case "FQDN":
			case "DSCP":
			case "10G6":
			case "MTU":
				ret.set(key, value === 1);
			default:
				ret.set(key, value);
				break;
		}
	}
	return ret;
}
