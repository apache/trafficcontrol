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

import { Component, OnInit, OnDestroy } from "@angular/core";
import { FormControl } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
// import { ColumnApi, GridApi, GridOptions, GridReadyEvent, RowNode } from "ag-grid-community";
import { ColumnApi, GridOptions, GridReadyEvent, RowNode } from "ag-grid-community";

import { Observable } from "rxjs";
import { map } from "rxjs/operators";

import { Interface, Server } from "../../../models/server";
import { ServerService } from "../../../services/api";
import { fuzzyScore, IPV4, serviceInterface } from "../../../utils";
import { SSHCellRendererComponent } from "../../table-components/ssh-cell-renderer/ssh-cell-renderer.component";

/**
 *
 */
interface AugmentedServer extends Server {
	/** The server's IPv4 service address */
	ipv4Address: string;
	/** The server's IPv6 service address */
	ipv6Address: string;
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

	/** Holds a reference to the GridOptions for the nested AG-Grid instance. */
	public gridOptions: GridOptions;
	/** Holds a reference to the nested grid's API. */
	// private gridAPI: GridApi | null = null;
	/** Holds a reference to the nested grid's column API. */
	private columnAPI: ColumnApi | null = null;

	/** All of the servers which should appear in the table. */
	public servers: Observable<Array<AugmentedServer>> | null = null;
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
			headerName: "IPv4 Gateway",
			hide: true,
			onCellClicked: null
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "ipv4Address",
			filter: true,
			headerName: "IPv4 Address",
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
			headerName: "IPv4 Subnet",
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
	public fuzzControl: FormControl = new FormControl("");

	// private userSubscription: Subscription;

	constructor(private readonly api: ServerService, private readonly route: ActivatedRoute, private readonly router: Router) {
		this.gridOptions = {
			doesExternalFilterPass: this.filter.bind(this),
			isExternalFilterPresent: this.shouldFilter.bind(this),
		};
		console.log("this.fuzzControl in constructor:", this.fuzzControl);
	}

	/** Initializes table data, loading it from Traffic Ops. */
	public ngOnInit(): void {
		this.servers = this.api.getServers().pipe(map(
			x => x.map(
					s => {
						const aug: AugmentedServer = {ipv4Address: "", ipv6Address: "", ...s};
						let inf: Interface;
						try {
							inf = serviceInterface(aug.interfaces);
						} catch (e) {
							console.error(`server #${s.id}:`, e);
							return aug;
						}
						for (const ip of inf.ipAddresses) {
							if (!ip.serviceAddress) {
								continue;
							}
							if (IPV4.test(ip.address)) {
								if (aug.ipv4Address !== "") {
									console.warn("found more than one IPv4 service address for server:", s.id);
								}
								aug.ipv4Address = ip.address;
							} else {
								if (aug.ipv6Address !== "") {
									console.warn("found more than one IPv6 service address for server:", s.id);
								}
								aug.ipv6Address = ip.address;
							}
						}
						return aug;
					}
				)
		));

		this.route.queryParamMap.subscribe(
			m => {
				const search = m.get("search");
				if (search) {
					this.fuzzControl.setValue(decodeURIComponent(search));
				}
			}
		);

		console.log("this.fuzzControl in ngOnInit:", this.fuzzControl);

	}

	/** Cleans up resources on destruction. */
	public ngOnDestroy(): void {
		// this.userSubscription.unsubscribe();
	}

	public setAPI(e: GridReadyEvent): void {
		this.columnAPI = e.columnApi;
		// this.gridAPI = e.api;
	}

	/**
	 * Lets AG-Grid know if it should be calling the external filter.
	 *
	 * @returns whether or not the external filter is to be used.
	 */
	public shouldFilter(): boolean {
		const should = true; //this.fuzzControl && this.fuzzControl.value !== "";
		console.log("should filter:", should);
		return should;
		// return this.fuzzControl && this.fuzzControl.value !== "";
	}

	/**
	 * Checks if a node passes the user's fuzzy search filter.
	 *
	 * @param node The row to check.
	 * @returns whether or not the row passes filtering
	 */
	public filter(node: RowNode): boolean {
		console.log("checking", node, "with", this.fuzzControl);
		console.log("this:", this);
		if (!this.columnAPI) {
			return true;
		}
		const visibleCols = new Set(this.columnAPI.getAllDisplayedColumns().map(x=>x.getColId()));
		for (const k in node.data) {
			if (!Object.prototype.hasOwnProperty.call(node.data, k) || !visibleCols.has(k)) {
				continue;
			}
			const value = node.data[k];
			let stringVal: string;
			switch (typeof value) {
				case "string":
					stringVal = value;
					break;
				case "number":
					stringVal = String(value);
					break;
				case "object":
					if (value instanceof Date) {
						stringVal = value.toLocaleString();
					} else {
						continue;
					}
					break;
				default:
					continue;
			}
			const score = fuzzyScore(stringVal.toLocaleLowerCase(), this.fuzzControl.value.toLocaleLowerCase());
			if ( score < Infinity) {
				console.log("server", node.data, "passed filter:", this.fuzzControl.value, "with score:", score, "on field:", k);
				return true;
			}
		}
		console.log("server", node.data, "failed filter:", this.fuzzControl.value);

		return false;
	}

	/** Update the URL's 'search' query parameter for the user's search input. */
	public updateURL(): void {
		if (this.gridOptions.api) {
			this.gridOptions.api.onFilterChanged();
		}
		if (this.fuzzControl.value === "") {
			this.router.navigate([], {queryParams: null, replaceUrl: true});
		} else if (this.fuzzControl.value) {
			this.router.navigate([], {queryParams: {search: this.fuzzControl.value}, replaceUrl: true});
		}
	}

}
