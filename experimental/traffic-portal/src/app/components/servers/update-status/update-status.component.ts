import { Component, EventEmitter, Input, OnInit, Output } from "@angular/core";
import { FormControl } from "@angular/forms";
import { merge, of } from "rxjs";
import { mergeAll } from "rxjs/operators";
import { Server, Status } from "src/app/models";
import { ServerService } from "src/app/services/api";

@Component({
	selector: "tp-update-status",
	styleUrls: ["./update-status.component.scss"],
	templateUrl: "./update-status.component.html"
})
export class UpdateStatusComponent implements OnInit {

	@Input() public open = false;
	@Input() public servers = new Array<Server>();
	@Output() public done = new EventEmitter<boolean>();

	public statuses = of(new Array<Status>());
	public currentStatus: null | number = null;
	public statusControl = new FormControl();
	public offlineReasonControl = new FormControl();

	public get isOffline(): boolean {
		const val = this.statusControl.value;
		return val !== null && val !== undefined && val.name !== "ONLINE" && val.name !== "REPORTED";
	}

	public get serverName(): string {
		const len = this.servers.length;
		if (len === 1) {
			return this.servers[0].hostName;
		}
		return `${len} servers`;
	}

	constructor(private readonly api: ServerService) { }

	/**
	 *
	 */
	public ngOnInit(): void {
		this.statuses = this.api.getStatuses();
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

	public submit(e: Event): void {
		e.preventDefault();
		e.stopPropagation();
		let observables;
		if (this.isOffline) {
			observables = this.servers.map(x=>this.api.updateStatus(x, this.statusControl.value.id, this.offlineReasonControl.value));
		} else {
			observables = this.servers.map(x=>this.api.updateStatus(x, this.statusControl.value.id));
		}
		merge(observables).pipe(mergeAll()).subscribe(
			() => {
				this.done.emit(true);
			},
			err => {
				console.error("something went wrong trying to update", this.serverName, "servers:", err);
				this.done.emit(false);
			}
		);
	}

	public cancel(): void {
		this.done.emit(false);
	}

}
