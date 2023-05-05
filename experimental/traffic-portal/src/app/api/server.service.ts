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
import type { RequestServer, RequestStatus, ResponseServer, ResponseStatus, Servercheck, ServerQueueResponse } from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * ServerService exposes API functionality related to Servers.
 */
@Injectable()
export class ServerService extends APIService {

	constructor(http: HttpClient) {
		super(http);
	}

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
	 * @returns The requested servers.
	 */
	public async getServers(): Promise<Array<ResponseServer>>;
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
}
