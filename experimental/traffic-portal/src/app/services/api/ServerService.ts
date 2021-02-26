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

import { Observable } from "rxjs";
import { map } from "rxjs/operators";

import { Server, Servercheck, Status } from "../../models";

import { APIService } from "./apiservice";

/**
 * Shared mapping for a single-Status request.
 *
 * @param s The raw response object from the API.
 * @returns The extracted and massaged Status object.
 */
function statusMap(s: [Status]): Status {
	s[0].lastUpdated = new Date((s[0].lastUpdated as unknown as string).replace("+00", "Z"));
	return s[0];
}

/**
 * ServerService exposes API functionality related to Servers.
 */
@Injectable({providedIn: "root"})
export class ServerService extends APIService {

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}

	public getServers(idOrName: number | string): Observable<Server>;
	public getServers(): Observable<Array<Server>>;
	/**
	 * Retrieves servers from the API.
	 *
	 * @param idOrName Specify either the integral, unique identifier (number) of a specific Server to retrieve, or its hostname (string).
	 * @returns An Observable that will emit the requested server(s).
	 */
	public getServers(idOrName?: number | string): Observable<Array<Server> | Server> {
		const path = "servers";
		if (idOrName !== undefined) {
			switch (typeof idOrName) {
				case "number":
					return this.get<[Server]>(path, undefined, {id: String(idOrName)}).pipe(map(
						r => {
							const srv = r[0];
							// These properties come in as strings
							if (srv.lastUpdated) {
								// Non-RFC3336 for some reason
								srv.lastUpdated = new Date((srv.lastUpdated as unknown as string).replace("+00", "Z"));
							}
							if (srv.statusLastUpdated) {
								srv.statusLastUpdated = new Date(srv.statusLastUpdated as unknown as string);
							}
							return srv;
						}
					));
				case "string":
					return this.get<Array<Server>>(path, undefined, {hostName: idOrName}).pipe(map(
						servers => {
							if (servers.length < 1) {
								throw new Error(`no such server '${idOrName}'`);
							}
							if (servers.length > 1) {
								console.warn(
									"Traffic Ops returned",
									servers.length,
									`servers with host name '${idOrName}' - selecting the first arbitrarily`
								);
							}
							const srv = servers[0];
							// These properties come in as strings
							if (srv.lastUpdated) {
								// Non-RFC3336 for some reason
								srv.lastUpdated = new Date((srv.lastUpdated as unknown as string).replace("+00", "Z"));
							}
							if (srv.statusLastUpdated) {
								srv.statusLastUpdated = new Date(srv.statusLastUpdated as unknown as string);
							}
							return srv;
						}
					));
			}
		}
		return this.get<Array<Server>>(path).pipe(map(
			r => r.map(
				s => {
					if (s.lastUpdated) {
						// Our dates are actually strings since JSON doesn't provide a native date type.
						// TODO: rework to use an interceptor
						const dateStr = (s.lastUpdated as unknown) as string;
						s.lastUpdated = new Date(dateStr.replace(" ", "T").replace(/\+00$/, "Z"));
						s.statusLastUpdated = s.statusLastUpdated ?
							new Date(s.statusLastUpdated as unknown as string) :
							s.statusLastUpdated;
					}
					return s;
				}
			)
		));
	}

	/**
	 * Creates a server.
	 *
	 * @param s The server to create.
	 * @returns An Observable that emits the server as created and returned by the API.
	 */
	public createServer(s: Server): Observable<Server> {
		return this.post<Server>("servers", s).pipe(map(
			srv => {
				srv.lastUpdated = new Date((srv.lastUpdated as unknown as string).replace("+00", "Z"));
				if (srv.statusLastUpdated) {
					srv.statusLastUpdated = new Date(srv.statusLastUpdated as unknown as string);
				}
				return srv;
			}
		));
	}

	public getServerChecks(): Observable<Servercheck[]>;
	public getServerChecks(id: number): Observable<Servercheck>;
	/**
	 * Fetches server "check" stats from Traffic Ops.
	 * Because the filter is not implemented on the server-side, the returned
	 * Observable<Servercheck> will throw an error if `id` does not exist.
	 *
	 * @param id If given, will return only the checks for the server with that ID.
	 * @todo Ideally this filter would be implemented server-side; the data set gets huge.
	 * @returns An observable that emits Serverchecks - or a single Servercheck if ID was given.
	 */
	public getServerChecks(id?: number): Observable<Servercheck | Servercheck[]> {
		const path = "servercheck";
		return this.get<Array<Servercheck>>(path).pipe(map(
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
		));
	}

	public getStatuses(idOrName: number | string): Observable<Status>;
	public getStatuses(): Observable<Array<Status>>;
	/**
	 * Retrieves Statuses from the API.
	 *
	 * @param idOrName An optional ID (number) or Name (string) used to fetch a single Status thereby identified.
	 * @returns An Observable that emits the requested Status(es).
	 */
	public getStatuses(idOrName?: number | string): Observable<Array<Status> | Status> {
		const path = "statuses";
		switch (typeof idOrName) {
			case "number":
				return this.get<[Status]>(path, {params: {id: String(idOrName)}}).pipe(map(statusMap));
			case "string":
				return this.get<[Status]>(path, {params: {name: idOrName}}).pipe(map(statusMap));
		}
		return this.get<Array<Status>>(path).pipe(map(
			ss => ss.map(
				s => {
					s.lastUpdated = new Date(s.lastUpdated as unknown as string);
					return s;
				}
			)
		));
	}

	/**
	 * Queues updates on a single server.
	 *
	 * @param server Either the server on which updates will be queued, or its integral, unique identifier.
	 * @returns The 'response' property of the TO server's response. See TO API docs.
	 */
	public queueUpdates(server: number | Server): Observable<{serverId: number; action: "queue"}> {
		let id: number;
		if (typeof server === "number") {
			id = server;
		} else if (!server.id) {
			throw new Error("server has no id");
		} else {
			id = server.id;
		}

		return this.post<{serverId: number; action: "queue"}>(`servers/${id}/queue_update`, {action: "queue"});
	}

	/**
	 * Clears updates on a single server.
	 *
	 * @param server Either the server for which updates will be cleared, or its integral, unique identifier.
	 * @returns The 'response' property of the TO server's response. See TO API docs.
	 */
	public clearUpdates(server: number | Server): Observable<{serverId: number; action: "dequeue"}> {
		let id: number;
		if (typeof server === "number") {
			id = server;
		} else if (!server.id) {
			throw new Error("server has no id");
		} else {
			id = server.id;
		}

		return this.post<{serverId: number; action: "dequeue"}>(`servers/${id}/queue_update`, {action: "dequeue"});
	}

	/**
	 * Updates a server's status.
	 *
	 * @param server Either the server that will have its status changed, or the integral, unique identifier thereof.
	 * @param status The name of the status to which to set the server.
	 * @param offlineReason The reason why the server was placed into a non-ONLINE or REPORTED status.
	 * @returns Nothing.
	 */
	public updateStatus(server: number | Server, status: string, offlineReason?: string): Observable<undefined> {
		let id: number;
		if (typeof server === "number") {
			id = server;
		} else if (!server.id) {
			throw new Error("server has no id");
		} else {
			id = server.id;
		}

		return this.put<undefined>(`servers/${id}/status`, {offlineReason, status});
	}
}
