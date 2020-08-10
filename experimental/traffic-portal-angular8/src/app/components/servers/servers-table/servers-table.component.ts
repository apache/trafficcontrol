import { Component, OnInit, OnDestroy } from "@angular/core";
import { FormControl } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";

import { Subscription } from "rxjs";

import { Server } from "../../../models/server";
import { ServerService } from "../../../services/api";
import { SSHCellRendererComponent } from "../../table-components/ssh-cell-renderer/ssh-cell-renderer.component";

// class SSHCellRenderer {
// 	private eGui: HTMLAnchorElement;
// 	private username: string;

// 	constructor (u: string) {
// 		this.username = u;
// 	}

// 	/** Input parameters used to set up the element that will be returned. */
// 	public init(params: {value: string}): void {
// 		this.eGui = document.createElement("A") as HTMLAnchorElement;
// 		this.eGui.href = `ssh://${this.username}@${params.value}`;
// 		this.eGui.setAttribute("target", "_blank");
// 		this.eGui.textContent = params.value;
// 	}

// 	/** Returns the HTML Element to put in the cell. */
// 	public getGui(): HTMLElement {
// 		return this.eGui;
// 	}
// }

class UpdateCellRenderer {
	private eGui: HTMLElement;

	/** Input parameters used to set up the element that will be returned. */
	public init(params: {value: boolean}): void {
		this.eGui = document.createElement("I");
		this.eGui.setAttribute("aria-hidden", "true");
		this.eGui.setAttribute("title", String(params.value));
		this.eGui.classList.add("fa", "fa-lg");
		if (params.value) {
			this.eGui.classList.add("fa-clock-o");
		} else {
			this.eGui.classList.add("fa-check");
		}
	}
}

/**
 * ServersTableComponent is the controller for the servers page - which
 * principally contains a table.
 */
@Component({
	selector: "servers-table",
	styleUrls: ["./servers-table.component.scss"],
	templateUrl: "./servers-table.component.html"
})
export class ServersTableComponent implements OnInit, OnDestroy {

	/** All of the servers which should appear in the table. */
	public servers: Array<Server>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			headerName: "Cache Group",
			field: "cachegroup",
			hide: false
		},
		{
			headerName: "CDN",
			field: "cdnName",
			hide: false
		},
		{
			headerName: "Domain",
			field: "domainName",
			hide: false
		},
		{
			headerName: "Host",
			field: "hostName",
			hide: false
		},
		{
			headerName: "HTTPS Port",
			field: "httpsPort",
			hide: true,
			filter: "agNumberColumnFilter"
		},
		{
			headerName: "Hash ID",
			field: "xmppId",
			hide: true
		},
		{
			headerName: "ID",
			field: "id",
			hide: true,
			filter: "agNumberColumnFilter"
		},
		{
			headerName: "ILO IP Address",
			field: "iloIpAddress",
			hide: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "ILO IP Gateway",
			field: "iloIpGateway",
			hide: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "ILO IP Netmask",
			field: "iloIpNetmask",
			hide: true
		},
		{
			headerName: "ILO Username",
			field: "iloUsername",
			hide: true
		},
		{
			headerName: "Interface Name",
			field: "interfaceName",
			hide: true
		},
		{
			headerName: "IPv6 Address",
			field: "ip6Address",
			hide: false
		},
		{
			headerName: "IPv6 Gateway",
			field: "ip6Gateway",
			hide: true
		},
		{
			headerName: "Last Updated",
			field: "lastUpdated",
			hide: true,
			filter: "agDateColumnFilter"
		},
		{
			headerName: "Mgmt IP Address",
			field: "mgmtIpAddress",
			hide: true
		},
		{
			headerName: "Mgmt IP Gateway",
			field: "mgmtIpGateway",
			hide: true,
			filter: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "Mgmt IP Netmask",
			field: "mgmtIpNetmask",
			hide: true,
			filter: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "Network Gateway",
			field: "ipGateway",
			hide: true,
			filter: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "Network IP",
			field: "ipAddress",
			hide: false,
			filter: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "Network MTU",
			field: "interfaceMtu",
			hide: true,
			filter: "agNumberColumnFilter"
		},
		{
			headerName: "Network Subnet",
			field: "ipNetmask",
			hide: true
		},
		{
			headerName: "Offline Reason",
			field: "offlineReason",
			hide: true
		},
		{
			headerName: "Phys Location",
			field: "physLocation",
			hide: true
		},
		{
			headerName: "Profile",
			field: "profile",
			hide: false
		},
		{
			headerName: "Rack",
			field: "rack",
			hide: true
		},
		{
			headerName: "Reval Pending",
			field: "revalPending",
			hide: true,
			filter: true,
			cellRenderer: "updateCellRenderer"
		},
		{
			headerName: "Router Hostname",
			field: "routerHostName",
			hide: true
		},
		{
			headerName: "Router Port Name",
			field: "routerPortName",
			hide: true
		},
		{
			headerName: "Status",
			field: "status",
			hide: false
		},
		{
			headerName: "TCP Port",
			field: "tcpPort",
			hide: true
		},
		{
			headerName: "Type",
			field: "type",
			hide: false
		},
		{
			headerName: "Update Pending",
			field: "updPending",
			hide: false,
			filter: true,
			cellRenderer: "updateCellRenderer"
		}
	];

	public components = {
		sshCellRenderer: SSHCellRendererComponent,
		// updateCellRenderer: new UpdateCellRenderer()
	};

	public get filteredServers(): Array<Server> {
		return this.servers.filter(x=>this.fuzzControl.value === "" || x.hostName.includes(this.fuzzControl.value));
	}

	/** Form controller for the user search input. */
	public fuzzControl: FormControl;

	// private userSubscription: Subscription;

	constructor (private readonly api: ServerService, private readonly route: ActivatedRoute, private readonly router: Router) {
		this.servers = [];
		this.fuzzControl = new FormControl("");
	}

	/** Initializes table data, loading it from Traffic Ops. */
	public ngOnInit(): void {
		this.api.getServers().subscribe(
			(r: Array<Server>) => {
				this.servers = r || [];
			}
		);

		// this.userSubscription = this.auth.currentUser.subscribe(
		// 	u => {
		// 		if (u && u.username) {
		// 			this.components.sshCellRenderer = new SSHCellRenderer(u.username);
		// 		}
		// 	}
		// );

		this.route.queryParamMap.subscribe(
			m => {
				const search = m.get("search");
				if (search) {
					this.fuzzControl.setValue(decodeURIComponent(search));
				}
			}
		);

	}

	/** Cleans up resources on destruction. */
	public ngOnDestroy(): void {
		// this.userSubscription.unsubscribe();
	}

	/** Update the URL's 'search' query parameter for the user's search input. */
	public updateURL(): void {
		if (this.fuzzControl.value === "") {
			this.router.navigate([], {replaceUrl: true, queryParams: null});
		} else if (this.fuzzControl.value) {
			this.router.navigate([], {replaceUrl: true, queryParams: {search: this.fuzzControl.value}});
		}
	}

}
