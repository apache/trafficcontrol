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

import { Injectable } from "@angular/core";
import type { RequestServer, RequestStatus, ResponseServer, ResponseStatus, Servercheck } from "trafficops-types";

import { CDNService, PhysicalLocationService, ProfileService, TypeService } from "..";

/**
 * Generates a `Servercheck` for a given `server`.
 *
 * @todo Inject the necessary services into the ServerService to be able to
 * generate this dynamically from IDs, instead of relying on optional names.
 *
 * @param server The server for which to generate a servercheck.
 * @returns A valid Servercheck for `server`.
 */
function serverCheck(server: ResponseServer): Servercheck {
	return {
		adminState: server.status ?? "SERVER HAD NO STATUS",
		cacheGroup: server.cachegroup ?? "SERVER HAD NO CACHE GROUP",
		hostName: server.hostName ?? "SERVER HAD NO HOST NAME",
		id: server.id,
		profile: server.profile ?? "SERVER HAD NO PROFILE",
		revalPending: server.revalPending,
		type: server.type ?? "SERVER HAD NO TYPE",
		updPending: server.updPending
	};
}

/**
 * ServerService exposes API functionality related to Servers.
 */
@Injectable()
export class ServerService {

	public servers = new Array<ResponseServer>();

	private readonly statuses: ResponseStatus[] = [
		{
			description: "Sever is administrative down and does not receive traffic.",
			id: 4,
			lastUpdated: new Date(),
			name: "ADMIN_DOWN"
		},
		{
			description: "Server is ignored by traffic router.",
			id: 5,
			lastUpdated: new Date(),
			name: "CCR_IGNORE"
		},
		{
			description: "Server is Offline. Not active in any configuration.",
			id: 1,
			lastUpdated: new Date(),
			name: "OFFLINE"
		},
		{
			description: "Server is online.",
			id: 2,
			lastUpdated: new Date(),
			name: "ONLINE"
		},
		{
			description: "Pre Production. Not active in any configuration.",
			id: 6,
			lastUpdated: new Date(),
			name: "PRE_PROD"
		},
		{
			description: "Server is online and reported in the health protocol.",
			id: 3,
			lastUpdated: new Date(),
			name: "REPORTED"
		}
	];

	private idCounter = 1;
	private statusIdCounter = 6;

	constructor(
		private readonly cdnService: CDNService,
		private readonly physLocService: PhysicalLocationService,
		private readonly typeService: TypeService,
		private readonly profileService: ProfileService
	){}

	public async getServers(idOrName: number | string): Promise<ResponseServer>;
	public async getServers(): Promise<Array<ResponseServer>>;
	/**
	 * Retrieves servers from the API.
	 *
	 * @param idOrName Specify either the integral, unique identifier (number) of a specific Server to retrieve, or its hostname (string).
	 * @returns The requested server(s).
	 */
	public async getServers(idOrName?: number | string): Promise<Array<ResponseServer> | ResponseServer> {
		if (idOrName !== undefined) {
			let server;
			switch (typeof idOrName) {
				case "number":
					server = this.servers.filter(s=>s.id === idOrName)[0];
					if (server === undefined) {
						throw new Error(`no such server: #${idOrName}`);
					}
					break;
				case "string":
					server = this.servers.filter(s=>s.hostName === idOrName)[0];
					if (server === undefined) {
						throw new Error(`no such server: '${idOrName}'`);
					}
					break;
			}
			return server;
		}
		return this.servers;
	}

	/**
	 * Creates a server.
	 *
	 * @param server The server to create.
	 * @returns The server as created and returned by the API.
	 */
	public async createServer(server: RequestServer): Promise<ResponseServer> {
		const cdn = await this.cdnService.getCDNs(server.cdnId);
		const physLoc = await this.physLocService.getPhysicalLocations(server.physLocationId);
		const profile = await this.profileService.getProfiles(server.profileId);
		const type = await this.typeService.getTypes(server.typeId);
		const status = await this.getStatuses(server.statusId);
		const newServer = {
			...server,
			// Due to circular dependency, name not resolved here
			cachegroup: "",
			cdnName: cdn.name,
			guid: server.guid ?? null,
			httpsPort: server.httpsPort ?? null,
			id: ++this.idCounter,
			iloIpAddress: server.iloIpAddress ?? null,
			iloIpGateway: server.iloIpGateway ?? null,
			iloIpNetmask: server.iloIpNetmask ?? null,
			iloPassword: server.iloPassword ?? null,
			iloUsername: server.iloUsername ?? null,
			lastUpdated: new Date(),
			mgmtIpAddress: server.mgmtIpAddress ?? null,
			mgmtIpGateway: server.mgmtIpGateway ?? null,
			mgmtIpNetmask: server.mgmtIpNetmask ?? null,
			offlineReason: server.offlineReason ?? null,
			physLocation: physLoc.name,
			profile: profile.name,
			profileDesc: profile.description,
			rack: server.rack ?? null,
			revalPending: false,
			routerHostName: server.routerHostName ?? null,
			routerPortName: server.routerPortName ?? null,
			status: status.name,
			statusLastUpdated: null,
			tcpPort: null,
			type: type.name,
			updPending: false,
			xmppId: ""
		};
		this.servers.push(newServer);
		return newServer;
	}

	public async getServerChecks(): Promise<Servercheck[]>;
	public async getServerChecks(id: number): Promise<Servercheck>;
	/**
	 * Fetches server "check" stats from Traffic Ops.
	 *
	 * @param id If given, will return only the checks for the server with that ID.
	 * @todo Ideally this filter would be implemented server-side; the data set gets huge.
	 * @returns Serverchecks - or a single Servercheck if ID was given.
	 */
	public async getServerChecks(id?: number): Promise<Servercheck | Servercheck[]> {
		if (id !== undefined) {
			const server = this.servers.filter(s=>s.id===id)[0];
			if (!server) {
				throw new Error(`no such server: #${id}`);
			}
			return serverCheck(server);
		}
		return this.servers.map(serverCheck);
	}

	public async getStatuses(idOrName: number | string): Promise<ResponseStatus>;
	public async getStatuses(): Promise<Array<ResponseStatus>>;
	/**
	 * Retrieves Statuses from the API.
	 *
	 * @param idOrName An optional ID (number) or Name (string) used to fetch a single Status thereby identified.
	 * @returns The requested Status(es).
	 */
	public async getStatuses(idOrName?: number | string): Promise<Array<ResponseStatus> | ResponseStatus> {
		if (idOrName !== undefined) {
			let status;
			if (typeof(idOrName) === "number") {
				status = this.statuses.filter(s=>s.id===idOrName)[0];
			} else {
				status = this.statuses.filter(s=>s.name===idOrName)[0];
			}
			if (!status) {
				throw new Error(`no such Status: ${idOrName}`);
			}
			return status;
		}
		return this.statuses;
	}

	/**
	 * Queues updates on a single server.
	 *
	 * @param server Either the server on which updates will be queued, or its integral, unique identifier.
	 * @returns The 'response' property of the TO server's response. See TO API docs.
	 */
	public async queueUpdates(server: number | ResponseServer): Promise<{serverId: number; action: "queue"}> {
		let id: number;
		if (typeof server === "number") {
			id = server;
		} else if (!server.id) {
			throw new Error("server has no id");
		} else {
			id = server.id;
		}

		const srv = this.servers.filter(s=>s.id===id)[0];
		if (!srv) {
			throw new Error(`no such Server: #${id}`);
		}

		srv.updPending = true;
		return {action: "queue", serverId: id};
	}

	/**
	 * Clears updates on a single server.
	 *
	 * @param server Either the server for which updates will be cleared, or its integral, unique identifier.
	 * @returns The 'response' property of the TO server's response. See TO API docs.
	 */
	public async clearUpdates(server: number | ResponseServer): Promise<{serverId: number; action: "dequeue"}> {
		let id: number;
		if (typeof server === "number") {
			id = server;
		} else if (!server.id) {
			throw new Error("server has no id");
		} else {
			id = server.id;
		}

		const srv = this.servers.filter(s=>s.id===id)[0];
		if (!srv) {
			throw new Error(`no such Server: #${id}`);
		}

		srv.updPending = false;
		return {action: "dequeue", serverId: id};
	}

	/**
	 * Updates a server's status.
	 *
	 * @param server Either the server that will have its status changed, or the integral, unique identifier thereof.
	 * @param statusName The name of the status to which to set the server.
	 * @param offlineReason The reason why the server was placed into a non-ONLINE or REPORTED status.
	 * @returns Nothing.
	 */
	public async updateStatus(server: number | ResponseServer, statusName: string, offlineReason?: string): Promise<void> {
		let id: number;
		if (typeof server === "number") {
			id = server;
		} else if (!server.id) {
			throw new Error("server has no id");
		} else {
			id = server.id;
		}

		const srv = this.servers.find(s=>s.id===id);
		if (!srv) {
			throw new Error(`no such Server: #${id}`);
		}

		const status = this.statuses.find(s=>s.name===statusName);
		if (!status) {
			throw new Error(`no such Status: '${statusName}'`);
		}
		if (status.id === undefined) {
			throw new Error(`Status with name '${statusName} has no ID`);
		}

		srv.status = statusName;
		srv.statusId = status.id;
		srv.offlineReason = offlineReason ?? null;
	}

	/**
	 * Creates a status.
	 *
	 * @param status The status details (name & description) to create. Description is an optional property in status.
	 * @returns The status as created and returned by the API.
	 */
	public async createStatus(status: RequestStatus): Promise<ResponseStatus> {
		const newStatus = {
			description: status.description ? status.description : null,
			id: ++this.statusIdCounter,
			lastUpdated: new Date(),
			name: status.name
		};
		this.statuses.push(newStatus);
		return newStatus;
	}

	/**
	 * Updates status Details.
	 *
	 * @param payload containes name and description for the status., unique identifier thereof.
	 * @param id The Status ID
	 */
	public async updateStatusDetail(payload: ResponseStatus, id: number): Promise<ResponseStatus> {
		const index = this.statuses.findIndex(u => u.id === id);
		if (index < 0) {
			throw new Error(`no such status with id: ${id}`);
		}
		const updated = {
			...payload,
			lastUpdated: new Date()
		} as { description: string; id: number; lastUpdated: Date; name: string };
		this.statuses[index] = updated;

		return updated;
	}

	/**
	 * Deletes a Status.
	 *
	 * @param id The ID of the Status to delete.
	 * @returns The deleted status.
	 */
	public async deleteStatus(id: number): Promise<ResponseStatus> {
		const idx = this.statuses.findIndex(j => j.id === id);
		if (idx < 0) {
			throw new Error(`no such status: #${id}`);
		}
		return this.statuses.splice(idx, 1)[0];
	}
}
