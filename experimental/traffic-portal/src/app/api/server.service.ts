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
import type { RequestServer, RequestStatus, ResponseServer, ResponseStatus, Servercheck } from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * ServerService exposes API functionality related to Servers.
 */
@Injectable()
export class ServerService extends APIService {

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}

	public async getServers(idOrName: number | string): Promise<ResponseServer>;
	public async getServers(): Promise<Array<ResponseServer>>;
	/**
	 * Retrieves servers from the API.
	 *
	 * @param idOrName Specify either the integral, unique identifier (number) of a specific Server to retrieve, or its hostname (string).
	 * @returns The requested server(s).
	 */
	public async getServers(idOrName?: number | string): Promise<Array<ResponseServer> | ResponseServer> {
		const path = "servers";
		if (idOrName !== undefined) {
			let servers;
			switch (typeof idOrName) {
				case "number":
					servers = await this.get<[ResponseServer]>(path, undefined, {id: String(idOrName)}).toPromise();
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
				console.warn(
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
		const path = "servercheck";
		return this.get<Array<Servercheck>>(path).toPromise().then(
			r => {
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
		);
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
		const path = "statuses";
		let ret;
		switch (typeof idOrName) {
			case "number":
				const response = await this.get<[ResponseStatus]>(path, undefined, { id: String(idOrName) }).toPromise();
				ret = response[0];
				break;
			case "string":
				ret = this.get<[ResponseStatus]>(path, {params: {name: idOrName}}).toPromise();
				break;
			default:
				ret = this.get<Array<ResponseStatus>>(path).toPromise();
		}
		return ret;
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

		return this.post<{serverId: number; action: "queue"}>(`servers/${id}/queue_update`, {action: "queue"}).toPromise();
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

		return this.post<{serverId: number; action: "dequeue"}>(`servers/${id}/queue_update`, {action: "dequeue"}).toPromise();
	}

	/**
	 * Updates a server's status.
	 *
	 * @param server Either the server that will have its status changed, or the integral, unique identifier thereof.
	 * @param status The name of the status to which to set the server.
	 * @param offlineReason The reason why the server was placed into a non-ONLINE or REPORTED status.
	 */
	public async updateStatus(server: number | ResponseServer, status: string, offlineReason?: string): Promise<undefined> {
		let id: number;
		if (typeof server === "number") {
			id = server;
		} else if (!server.id) {
			throw new Error("server has no id");
		} else {
			id = server.id;
		}

		return this.put(`servers/${id}/status`, {offlineReason, status}).toPromise();
	}

	/**
	 * Creating new Status.
	 *
	 * @param data containes name and description for the status.
	 * @returns The 'response' property of the TO status response. See TO API docs.
	 */
	public async createStatus(payload: RequestStatus): Promise<ResponseStatus> {
		return this.post<ResponseStatus>("statuses", payload).toPromise();
	}

	/**
	 * Updates status Details.
	 *
	 * @param data containes name and description for the status., unique identifier thereof.
	 * @param id The Status ID
	 */
	public async updateStatusDetail(payload: ResponseStatus, id: number): Promise<ResponseStatus> {
		return this.put<ResponseStatus>(`statuses/${id}`, payload).toPromise();
	}

	/**
	 * Deletes an existing Status.
	 *
	 * @param id The Status ID
	 */
	public async deleteStatus(id: number): Promise<void> {
		return this.delete(`statuses/${id}`).toPromise();
	}
}
