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
import {Component, Inject, type OnInit} from "@angular/core";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import type { ResponseServer, ResponseStatus } from "trafficops-types";

import { ServerService } from "src/app/api/server.service";
import { LoggingService } from "src/app/shared/logging.service";

/**
 * UpdateStatusComponent is the controller for the "Update Server Status" dialog box.
 */
@Component({
	selector: "tp-update-status[servers]",
	styleUrls: ["./update-status.component.scss"],
	templateUrl: "./update-status.component.html"
})
export class UpdateStatusComponent implements OnInit {

	/** The possible statuses of a server. */
	public statuses = new Array<ResponseStatus>();
	/** The ID of the current status of the server, or null if the servers have disparate statuses. */
	public currentStatus: null | number = null;

	public status: ResponseStatus | null = null;

	public servers: Array<ResponseServer>;

	public offlineReason = "";

	/** Tells whether the user's selected status is considered "OFFLINE". */
	public get isOffline(): boolean {
		return this.status !== null && this.status !== undefined &&
			this.status.name !== "ONLINE" && this.status.name !== "REPORTED";
	}

	/** An appropriate title for the server or collection of servers being updated. */
	public get serverName(): string {
		const len = this.servers.length;
		if (len === 1) {
			return this.servers[0].hostName;
		}
		return `${len} servers`;
	}

	constructor(
		private readonly dialogRef: MatDialogRef<UpdateStatusComponent>,
		@Inject(MAT_DIALOG_DATA) private readonly dialogServers: Array<ResponseServer>,
		private readonly api: ServerService,
		private readonly log: LoggingService,
	) {
		this.servers = this.dialogServers;
	}

	/**
	 * Sets up the necessary data to complete the form.
	 */
	public ngOnInit(): void {
		this.api.getStatuses().then(
			ss => {
				this.statuses = ss;
			}
		).catch(
			e => {
				this.log.error("Failed to get Statuses:", e);
			}
		);
		if (this.servers.length < 1) {
			return;
		}

		if (this.servers.length === 1) {
			this.currentStatus = this.servers[0].statusId;
		} else {
			const stat = this.servers[0].statusId;
			if (this.servers.every(x=>x.statusId === stat)) {
				this.currentStatus = stat;
			} else  {
				this.currentStatus = null;
			}
		}
	}

	/**
	 * Triggered when the user submits the form; this attempts to update the server(s) and emits
	 * from the `done` Output a value that depends on success or failure of the request.
	 *
	 * @param e The submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		let observables;
		if (this.isOffline) {
			observables = this.servers.map(
				async x=> this.api.updateStatus(x, this.status?.name ?? "", this.offlineReason)
			);
		} else {
			observables = this.servers.map(async x=>this.api.updateStatus(x, this.status?.name ?? ""));
		}
		try {
			await Promise.all(observables);
			this.dialogRef.close(true);
		} catch (err) {
			this.log.error("something went wrong trying to update", this.serverName, "servers:", err);
			this.dialogRef.close(false);
		}
	}

	/**
	 * Emits from the `done` Output indicating the action was cancelled.
	 */
	public cancel(): void {
		this.dialogRef.close(false);
	}

}
