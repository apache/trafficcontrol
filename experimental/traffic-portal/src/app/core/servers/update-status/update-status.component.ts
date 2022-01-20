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
import { Component, EventEmitter, HostListener, Input, type OnInit, Output } from "@angular/core";
import { FormControl } from "@angular/forms";

import { ServerService } from "src/app/api/server.service";
import type { Server, Status } from "src/app/models";

/**
 * UpdateStatusComponent is the controller for the "Update Server Status" dialog box.
 */
@Component({
	selector: "tp-update-status[servers]",
	styleUrls: ["./update-status.component.scss"],
	templateUrl: "./update-status.component.html"
})
export class UpdateStatusComponent implements OnInit {

	/** The servers being updated. */
	@Input() public servers = new Array<Server>();
	/** Emits 'false' when the status update is cancelled (or fails), 'true' when it completes successfully. */
	@Output() public done = new EventEmitter<boolean>();

	/**
	 * Captures keypresses and emits 'false' from the 'done' output if the user presses the escape key.
	 *
	 * @param e The captured 'keydown' event.
	 */
	@HostListener("document:keydown", ["$event"]) public closeOnEscape(e: KeyboardEvent): void {
		if (e.code === "Escape" || e.code === "Esc") {
			this.done.emit(false);
		}
	}

	/** The possible statuses of a server. */
	public statuses = new Array<Status>();
	/** The ID of the current status of the server, or null if the servers have disparate statuses. */
	public currentStatus: null | number = null;

	/** Form control for the new status selection. */
	public statusControl = new FormControl();
	/** Form control for the offline reason input. */
	public offlineReasonControl = new FormControl();

	/** Tells whether the user's selected status is considered "OFFLINE". */
	public get isOffline(): boolean {
		const val = this.statusControl.value;
		return val !== null && val !== undefined && val.name !== "ONLINE" && val.name !== "REPORTED";
	}

	/** An appropriate title for the server or collection of servers being updated. */
	public get serverName(): string {
		const len = this.servers.length;
		if (len === 1) {
			return this.servers[0].hostName;
		}
		return `${len} servers`;
	}

	/** Constructor. */
	constructor(private readonly api: ServerService) { }

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
				console.error("Failed to get Statuses:", e);
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
	public submit(e: Event): void {
		e.preventDefault();
		e.stopPropagation();
		let observables;
		if (this.isOffline) {
			observables = this.servers.map(async x=>this.api.updateStatus(x, this.statusControl.value.id, this.offlineReasonControl.value));
		} else {
			observables = this.servers.map(async x=>this.api.updateStatus(x, this.statusControl.value.id));
		}
		Promise.all(observables).then(
			() => {
				this.done.emit(true);
			},
			err => {
				console.error("something went wrong trying to update", this.serverName, "servers:", err);
				this.done.emit(false);
			}
		);
	}

	/**
	 * Emits from the `done` Output indicating the action was cancelled.
	 */
	public cancel(): void {
		this.done.emit(false);
	}

}
