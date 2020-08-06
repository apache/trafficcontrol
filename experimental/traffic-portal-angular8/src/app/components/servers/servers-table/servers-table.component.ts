import { Component, OnInit } from "@angular/core";

import { Server } from "../../../models/server";
import { ServerService } from "../../../services/api";

/**
 * ServersTableComponent is the controller for the servers page - which
 * principally contains a table.
 */
@Component({
	selector: "servers-table",
	styleUrls: ["./servers-table.component.scss"],
	templateUrl: "./servers-table.component.html"
})
export class ServersTableComponent implements OnInit {

	/** All of the servers which should appear in the table. */
	public servers: Array<Server>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{headerName: "ID", field: "id" },
		{headerName: "Hostname", field: "hostName" },
		{headerName: "Profile", field: "profile"}
	];

	constructor(private readonly api: ServerService) {
		this.servers = [];
	}

	/** Initializes table data, loading it from Traffic Ops. */
	public ngOnInit(): void {
		this.api.getServers().subscribe(
			(r: Array<Server>) => {
				this.servers = r;
			}
		);
	}

}
