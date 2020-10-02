import { Component, OnInit, OnDestroy } from "@angular/core";
import { FormControl } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { Observable } from "rxjs";


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

interface UpdateCellRendererOptions {
	/** the value being rendered */
	value: boolean;
}

class UpdateCellRenderer {
	private eGui: HTMLElement;

	/** Input parameters used to set up the element that will be returned. */
	public init(params: UpdateCellRendererOptions): void {
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
	public servers: Observable<Array<Server>>;
	// public servers: Array<Server>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "cachegroup",
			headerName: "Cache Group",
			hide: false
		},
		{
			field: "cdnName",
			headerName: "CDN",
			hide: false
		},
		{
			field: "domainName",
			headerName: "Domain",
			hide: false
		},
		{
			field: "hostName",
			headerName: "Host",
			hide: false
		},
		{
			field: "httpsPort",
			filter: "agNumberColumnFilter",
			headerName: "HTTPS Port",
			hide: true,
		},
		{
			field: "xmppId",
			headerName: "Hash ID",
			hide: true
		},
		{
			field: "id",
			filter: "agNumberColumnFilter",
			headerName: "ID",
			hide: true,
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "iloIpAddress",
			headerName: "ILO IP Address",
			hide: true,
			onCellClicked: null
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "iloIpGateway",
			headerName: "ILO IP Gateway",
			hide: true,
			onCellClicked: null
		},
		{
			field: "iloIpNetmask",
			headerName: "ILO IP Netmask",
			hide: true
		},
		{
			field: "iloUsername",
			headerName: "ILO Username",
			hide: true
		},
		{
			field: "interfaceName",
			headerName: "Interface Name",
			hide: true
		},
		{
			field: "ip6Address",
			headerName: "IPv6 Address",
			hide: false
		},
		{
			field: "ip6Gateway",
			headerName: "IPv6 Gateway",
			hide: true
		},
		{
			field: "lastUpdated",
			filter: "agDateColumnFilter",
			headerName: "Last Updated",
			hide: true,
		},
		{
			field: "mgmtIpAddress",
			headerName: "Mgmt IP Address",
			hide: true
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "mgmtIpGateway",
			filter: true,
			headerName: "Mgmt IP Gateway",
			hide: true,
			onCellClicked: null
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "mgmtIpNetmask",
			filter: true,
			headerName: "Mgmt IP Netmask",
			hide: true,
			onCellClicked: null
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "ipGateway",
			filter: true,
			headerName: "Network Gateway",
			hide: true,
			onCellClicked: null
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "ipAddress",
			filter: true,
			headerName: "Network IP",
			hide: false,
			onCellClicked: null
		},
		{
			field: "interfaceMtu",
			filter: "agNumberColumnFilter",
			headerName: "Network MTU",
			hide: true,
		},
		{
			field: "ipNetmask",
			headerName: "Network Subnet",
			hide: true
		},
		{
			field: "offlineReason",
			headerName: "Offline Reason",
			hide: true
		},
		{
			field: "physLocation",
			headerName: "Phys Location",
			hide: true
		},
		{
			field: "profile",
			headerName: "Profile",
			hide: false
		},
		{
			field: "rack",
			headerName: "Rack",
			hide: true
		},
		{
			cellRenderer: "updateCellRenderer",
			field: "revalPending",
			filter: true,
			headerName: "Reval Pending",
			hide: true,
		},
		{
			field: "routerHostName",
			headerName: "Router Hostname",
			hide: true
		},
		{
			field: "routerPortName",
			headerName: "Router Port Name",
			hide: true
		},
		{
			field: "status",
			headerName: "Status",
			hide: false
		},
		{
			field: "tcpPort",
			headerName: "TCP Port",
			hide: true
		},
		{
			field: "type",
			headerName: "Type",
			hide: false
		},
		{
			cellRenderer: "updateCellRenderer",
			field: "updPending",
			filter: true,
			headerName: "Update Pending",
			hide: false,
		}
	];

	/** Passed as components to the ag-grid API */
	public components = {
		sshCellRenderer: SSHCellRendererComponent,
		// updateCellRenderer: new UpdateCellRenderer()
	};

	/** a list of all servers that match the current filter */
	// public get filteredServers(): Array<Server> {
	// 	return this.servers.filter(x=>this.fuzzControl.value === "" || x.hostName.includes(this.fuzzControl.value));
	// }

	/** Form controller for the user search input. */
	public fuzzControl: FormControl;

	// private userSubscription: Subscription;

	constructor (private readonly api: ServerService, private readonly route: ActivatedRoute, private readonly router: Router) {
		// this.servers = [];
		this.fuzzControl = new FormControl("");
	}

	/** Initializes table data, loading it from Traffic Ops. */
	public ngOnInit(): void {
		this.servers = this.api.getServers();
		// this.api.getServers().subscribe(
		// 	(r: Array<Server>) => {
		// 		this.servers = r || [];
		// 	}
		// );

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
