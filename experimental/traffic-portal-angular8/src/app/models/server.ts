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

export interface Server {
	cachegroup?: string;
	cachegroupId: number;
	cdnId: number;
	cdnName?: string;
	domainName: string;
	hostName: string;
	httpsPort: number;
	id?: number;
	iloIpAddress?: string;
	iloIpGateWay?: string;
	iloIPNetmask?: string;
	iloPassword?: string;
	iloUsername?: string;
	ip6Address?: string;
	ip6Gateway?: string;
	ipAddress: string;
	ipGateway?: string;
	ipNetmask?: string;
	lastUpdated?: Date;
	mgmtIpAddress?: string;
	mgmtIpGateway?: string;
	mgmtIpNetmask?: string;
	offlineReason?: string;
	physLocation?: string;
	physLocationId: number;
	profile?: string;
	profileDesc?: string;
	profileId: number;
	rack?: string;
	revalPending: boolean;
	routerHostName?: string;
	routerPortName?: string;
	status?: string;
	statusId: number;
	tcpPort?: number;
	type?: string;
	typeId: number;
	updPending: boolean;
	xmppId?: string;
	xmppPasswd?: string;
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
