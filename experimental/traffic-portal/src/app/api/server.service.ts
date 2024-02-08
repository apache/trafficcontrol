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

import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import type {
	Interface,
	IPAddress,
	RequestServer,
	RequestServerCapability,
	RequestStatus,
	ResponseServer,
	ResponseServerCapability,
	ResponseStatus,
	Server,
	ServerCapability,
	Servercheck,
	ServerQueueResponse,
} from "trafficops-types";

import { LoggingService } from "../shared/logging.service";

import { APIService } from "./base-api.service";

/**
 * ServerService exposes API functionality related to Servers.
 */
@Injectable()
export class ServerService extends APIService {

	constructor(http: HttpClient, private readonly log: LoggingService) {
		super(http);
	}

	/**
	 * Retrieves servers from the API.
	 *
	 * @returns The requested servers.
	 */
	public async getServers(): Promise<Array<ResponseServer>>;
	/**
	 * Retrieves a single server from the API.
	 *
	 * @param idOrName Either the integral, unique identifier (number) of a
	 * specific Server to retrieve, or its hostname (string).
	 * @returns The requested server. Note that hostNames are **not** unique,
	 * despite the vast number of places ATC components assume that they are. In
	 * the event that more than one server shares a hostName, this method will
	 * arbitrarily choose the first and log a warning, rather than throwing an
	 * error.
	 */
	public async getServers(idOrName: number | string): Promise<ResponseServer>;
	/**
	 * Retrieves servers from the API.
	 *
	 * @param idOrName Optionally specify either the integral, unique identifier
	 * (number) of a specific Server to retrieve, or its hostname (string).
	 * @returns The requested server(s).
	 */
	public async getServers(idOrName?: number | string): Promise<Array<ResponseServer> | ResponseServer> {
		const path = "servers";
		if (idOrName !== undefined) {
			let servers;
			switch (typeof idOrName) {
				case "number":
					servers = await this.get<[ResponseServer]>(path, undefined, {id: idOrName}).toPromise();
					break;
				case "string":
					servers = await this.get<Array<ResponseServer>>(path, undefined, {hostName: idOrName}).toPromise();
			}
			if (servers.length < 1) {
				throw new Error(`no such server '${idOrName}'`);
			}
			// This is, unfortunately, possible, despite the many assumptions to
			// the contrary.
			if (servers.length > 1) {
				this.log.warn(
					"Traffic Ops returned",
					servers.length,
					`servers with host name '${idOrName}' - selecting the first arbitrarily`
				);
			}
			return servers[0];
		}
		return this.get<Array<ResponseServer>>(path).toPromise();
	}

	/**
	 * Creates a server.
	 *
	 * @param s The server to create.
	 * @returns The server as created and returned by the API.
	 */
	public async createServer(s: RequestServer): Promise<ResponseServer> {
		return this.post<ResponseServer>("servers", s).toPromise();
	}

	/**
	 * Updates a server by the given payload
	 *
	 * @param serverOrID The server object or id to be deleted
	 * @param payload The server payload to update with.
	 */
	public async updateServer(serverOrID: ResponseServer | number, payload?: RequestServer): Promise<ResponseServer> {
		let id;
		let body;
		if (typeof(serverOrID) === "number") {
			if(!payload) {
				throw new TypeError("invalid call signature - missing request paylaod");
			}
			body = payload;
			id = +serverOrID;
		} else {
			body = serverOrID;
			id = serverOrID.id;
		}

		return this.put<ResponseServer>(`servers/${id}`, body).toPromise();
	}

	/**
	 * Fetches server "check" stats from Traffic Ops.
	 *
	 * @returns All Serverchecks Traffic Ops has.
	 */
	public async getServerChecks(): Promise<Servercheck[]>;
	/**
	 * Fetches a server's "check" stats from Traffic Ops.
	 *
	 * @param id The ID of the server whose "checks" will be returned.
	 * @returns The Servercheck for the server identified by `id`.
	 */
	public async getServerChecks(id: number): Promise<Servercheck>;
	/**
	 * Fetches server "check" stats from Traffic Ops.
	 *
	 * @param id If given, will return only the checks for the server with that
	 * ID.
	 * @todo Ideally this filter would be implemented server-side; the data set
	 * gets huge.
	 * @returns Serverchecks - or a single Servercheck if ID was given.
	 */
	public async getServerChecks(id?: number): Promise<Servercheck | Servercheck[]> {
		const path = "servercheck";
		const r = await this.get<Array<Servercheck>>(path).toPromise();
		if (id) {
			for (const sc of r) {
				if (sc.id === id) {
					return sc;
				}
			}
			throw new Error(`no server #${id} found in checks response`);
		}
		return r;
	}

	/**
	 * Retrieves a specific Status from the API.
	 *
	 * @param idOrName The ID (number) or Name (string) of a single Status to be
	 * retrieved.
	 * @returns The requested Status.
	 */
	public async getStatuses(idOrName: number | string): Promise<ResponseStatus>;
	/**
	 * Retrieves Statuses from the API.
	 *
	 * @returns The requested Statuses.
	 */
	public async getStatuses(): Promise<Array<ResponseStatus>>;
	/**
	 * Retrieves Statuses from the API.
	 *
	 * @param idOrName An optional ID (number) or Name (string) used to fetch a
	 * single Status thereby identified.
	 * @returns The requested Status(es).
	 */
	public async getStatuses(idOrName?: number | string): Promise<Array<ResponseStatus> | ResponseStatus> {
		const path = "statuses";
		if (idOrName !== undefined) {
			let params;
			if (typeof(idOrName) === "number") {
				params = {id: idOrName};
			} else {
				params = {name: idOrName};
			}
			const ret = await this.get<[ResponseStatus]>(path, undefined, params).toPromise();
			if (ret.length !== 1) {
				throw new Error(`Traffic Ops reported ${ret.length} Statuses by identifier '${idOrName}'`);
			}
			return ret[0];
		}
		return this.get<Array<ResponseStatus>>(path).toPromise();
	}

	/**
	 * Queues updates on a single server.
	 *
	 * @param server Either the server on which updates will be queued, or its
	 * integral, unique identifier.
	 * @returns The 'response' property of the TO server's response. See TO API
	 * docs.
	 */
	public async queueUpdates(server: number | ResponseServer): Promise<ServerQueueResponse> {
		const id = typeof(server) === "number" ? server : server.id;
		return this.post<ServerQueueResponse>(`servers/${id}/queue_update`, {action: "queue"}).toPromise();
	}

	/**
	 * Clears updates on a single server.
	 *
	 * @param server Either the server for which updates will be cleared, or its
	 * integral, unique identifier.
	 * @returns The 'response' property of the TO server's response. See TO API
	 * docs.
	 */
	public async clearUpdates(server: number | ResponseServer): Promise<ServerQueueResponse> {
		const id = typeof(server) === "number" ? server : server.id;
		return this.post<ServerQueueResponse>(`servers/${id}/queue_update`, {action: "dequeue"}).toPromise();
	}

	/**
	 * Updates a server's status.
	 *
	 * @param server Either the server that will have its status changed, or the
	 * integral, unique identifier thereof.
	 * @param newStatus Either the status to which to set the server, or the
	 * name thereof.
	 * @param offlineReason The reason why the server was placed into a
	 * non-ONLINE or REPORTED status.
	 */
	public async updateStatus(server: number | ResponseServer, newStatus: string | ResponseStatus, offlineReason?: string): Promise<void> {
		const id = typeof(server) === "number" ? server : server.id;
		const status = typeof(newStatus) === "string" ? newStatus : newStatus.name;
		return this.put(`servers/${id}/status`, {offlineReason, status}).toPromise();
	}

	/**
	 * Creating new Status.
	 *
	 * @param status The status to create.
	 * @returns The created status.
	 */
	public async createStatus(status: RequestStatus): Promise<ResponseStatus> {
		return this.post<ResponseStatus>("statuses", status).toPromise();
	}

	/**
	 * Updates status Details.
	 *
	 * @param status The status to update.
	 * @returns The updated status.
	 */
	public async updateStatusDetail(status: ResponseStatus): Promise<ResponseStatus> {
		return this.put<ResponseStatus>(`statuses/${status.id}`, status).toPromise();
	}

	/**
	 * Deletes an existing Status.
	 *
	 * @param statusId The Status ID
	 */
	public async deleteStatus(statusId: number | ResponseStatus): Promise<ResponseStatus> {
		const id = typeof (statusId) === "number" ? statusId : statusId.id;
		return this.delete<ResponseStatus>(`statuses/${id}`).toPromise();
	}

	/**
	 * Deletes an existing server.
	 *
	 * @param server The Server to be deleted, or just its ID.
	 * @returns The deleted server.
	 */
	public async deleteServer(server: number | ResponseServer): Promise<ResponseServer> {
		const id =  typeof(server) === "number" ? server : server.id;
		const path = `servers/${id}`;
		return this.delete<ResponseServer>(path).toPromise();

	}

	/**
	 * Retrieves Server Capabilities from Traffic Ops.
	 *
	 * @returns All requested Capabilities.
	 */
	public async getCapabilities(): Promise<Array<ResponseServerCapability>>;
	/**
	 * Retrieves a specific Server Capability from Traffic Ops.
	 *
	 * @param name The name of the requested Server Capability.
	 * @returns The requested Capability.
	 * @throws {Error} if Traffic Ops responds with any number of Capabilities
	 * besides exactly one.
	 */
	public async getCapabilities(name: string): Promise<ResponseServerCapability>;
	/**
	 * Retrieves one or more Server Capabilities from Traffic Ops.
	 *
	 * @param name If given, only the Capability with this name will be
	 * returned.
	 * @returns Any and all requested Capabilities.
	 * @throws {Error} if a Capability is requested by name, but Traffic Ops
	 * responds with any number of Capabilities besides exactly one.
	 */
	public async getCapabilities(name?: string): Promise<Array<ResponseServerCapability> | ResponseServerCapability> {
		const path = "server_capabilities";
		if (name) {
			const resp = await this.get<[ResponseServerCapability]>(path, undefined, {name}).toPromise();
			if (resp.length !== 1) {
				throw new Error(`Traffic Ops responded with ${resp.length} Capabilities with name '${name}'`);
			}
			return resp[0];
		}
		return this.get<Array<ResponseServerCapability>>(path).toPromise();
	}

	/**
	 * Deletes a Server Capability.
	 *
	 * @param cap The Capability to be deleted, or just its name.
	 */
	public async deleteCapability(cap: string | ServerCapability): Promise<void> {
		const name = typeof(cap) === "string" ? cap : cap.name;
		return this.delete("server_capabilities", undefined, {name}).toPromise();
	}

	/**
	 * Replaces an existing Server Capability definition with a new one.
	 *
	 * @param name The Capability's current Name.
	 * @param cap The Capability with desired modifications made.
	 * @returns The modified Capability.
	 */
	public async updateCapability(name: string, cap: ServerCapability): Promise<ResponseServerCapability> {
		return this.put<ResponseServerCapability>("server_capabilities", cap, {name}).toPromise();
	}

	/**
	 * Creates a new Server Capability.
	 *
	 * @param cap The new Capability.
	 * @returns The created Capability.
	 */
	public async createCapability(cap: RequestServerCapability): Promise<ResponseServerCapability> {
		return this.post<ResponseServerCapability>("server_capabilities", cap).toPromise();
	}

	/**
	 * Gets the "service" interface for a server; that is, the interface that
	 * contains service addresses.
	 *
	 * @param server Either the server for which to find the "service"
	 * interface, or just the interfaces thereof.
	 * @returns The network interface that contains the service addresses.
	 * @throws {Error} If no service addresses are found on any interface.
	 */
	public static getServiceInterface(server: Server | Interface[]): Interface {
		const infs = Array.isArray(server) ? server : server.interfaces;
		for (const inf of infs) {
			for (const addr of inf.ipAddresses) {
				if (addr.serviceAddress) {
					return inf;
				}
			}
		}
		throw new Error("no service addresses found");
	}

	/**
	 * Pulls apart an IP address with a CIDR-notation suffix into a plain
	 * address (with no suffix) and a netmask that represents the same subnet
	 * as the CIDR-notation suffix.
	 *
	 * @param addr The address from which to extract the netmask.
	 * @returns The address without a netmask and the netmask itself (if one
	 * could be found; otherwise it'll be `undefined`).
	 */
	public static extractNetmask(addr: IPAddress | string): [string, string | undefined] {
		let addrStr = typeof(addr) === "string" ? addr : addr.address;
		let maskStr;
		if (addrStr.includes("/")) {
			const parts = addrStr.split("/");
			addrStr = parts[0];
			let masklen = Number(parts[1]);

			const mask = [];
			for (let k = 0; k < 4; ++k) {
				const n = Math.min(masklen, 8);
				mask.push(256 - Math.pow(2, 8 - n));
				masklen -= n;
			}
			maskStr = mask.join(".");
		}
		return [addrStr, maskStr];
	}
}
