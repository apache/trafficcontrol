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

import { Component , OnInit} from "@angular/core";
import { FormControl } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import type { ITooltipParams } from "ag-grid-community";
import { BehaviorSubject } from "rxjs";
import { ResponseCDN, ResponseServer, serviceAddresses } from "trafficops-types";

import { CDNService, ServerService } from "src/app/api";
import { UpdateStatusComponent } from "src/app/core/servers/update-status/update-status.component";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import {
	CollectionChoiceDialogComponent, CollectionChoiceDialogData
} from "src/app/shared/dialogs/collection-choice-dialog/collection-choice-dialog.component";
import type {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink,
	TableTitleButton
} from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * AugmentedServer has fields that give direct access to its service addresses without needing to recalculate them.
 */
export interface AugmentedServer extends ResponseServer {
	/** The server's IPv4 service address */
	ipv4Address: string;
	/** The server's IPv6 service address */
	ipv6Address: string;
}

/**
 * Converts a server to an "augmented" server.
 *
 * @param s The server to convert.
 * @returns The converted server.
 */
export function augment(s: ResponseServer): AugmentedServer {
	const [ipv4Address, ipv6Address] = serviceAddresses(s.interfaces);
	return {
		...s,
		ipv4Address: ipv4Address ? ipv4Address.address : "" ,
		ipv6Address: ipv6Address ? ipv6Address.address : ""
	};
}

/**
 * Checks if a server is a Cache Server.
 *
 * @param data The server to check.
 * @returns Whether or not 'data' is a Cache Server.
 */
export function serverIsCache(data: AugmentedServer): boolean {
	if (!data || !data.type) {
		return false;
	}
	return data.type.startsWith("EDGE") || data.type.startsWith("MID");
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
export class ServersTableComponent implements OnInit {

	/** All of the servers which should appear in the table. */
	public servers: Promise<Array<AugmentedServer>> | null = null;

	/** All of the CDNs (on which a user might (de/)queue updates). */
	public readonly cdns: Promise<Array<ResponseCDN>>;

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
			onCellClicked: undefined
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "iloIpGateway",
			headerName: "ILO IP Gateway",
			hide: true,
			onCellClicked: undefined
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
			onCellClicked: undefined
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "mgmtIpNetmask",
			filter: true,
			headerName: "Mgmt IP Netmask",
			hide: true,
			onCellClicked: undefined
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "ipGateway",
			filter: true,
			headerName: "IPv4 Gateway",
			hide: true,
			onCellClicked: undefined
		},
		{
			cellRenderer: "sshCellRenderer",
			field: "ipv4Address",
			filter: true,
			headerName: "IPv4 Address",
			hide: false,
			onCellClicked: undefined
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
			hide: true,
			valueFormatter: ({data}: {data: AugmentedServer}): string => `${data.physLocation} (#${data.physLocationId})`
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
			filter: "tpBooleanFilter",
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
			hide: false,
			tooltipValueGetter(params: ITooltipParams): string {
				if (!params.value || params.value === "ONLINE" || params.value === "REPORTED") {
					return "";
				}
				return `${params.value}: ${params.data.offlineReason}`;
			}
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
			filter: "tpBooleanFilter",
			headerName: "Update Pending",
			hide: false,
		}
	];

	/** Definitions for the more menu buttons */
	public moreMenuButtons: Array<TableTitleButton> = [
		{
			action: "queue",
			text: "Queue Server Updates"
		},
		{
			action: "dequeue",
			text: "Clear Server Updates"
		}
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<AugmentedServer> = {
		href: (row: AugmentedServer): string => `/core/servers/${row.id}`
	};

	/** Definitions for the context menu items (which act on augmented server data). */
	public contextMenuItems: Array<ContextMenuItem<AugmentedServer>> = [
		{
			href: (row: AugmentedServer): string => `${row.id}`,
			name: "View Server Details"
		},
		{
			href: (row: AugmentedServer): string => `${row.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			href: (row: AugmentedServer): string => `/core/cache-groups/${row.cachegroupId}`,
			name: "View Cache Group"
		},
		{
			href: (row: AugmentedServer): string => `/core/phys-locs/${row.physLocationId}`,
			name: "View Physical Location"
		},
		{
			action: "updateStatus",
			multiRow: true,
			name: "Update Status"
		},
		{
			action: "queue",
			disabled: (data: Array<AugmentedServer>): boolean => !data.every(serverIsCache),
			multiRow: true,
			name: "Queue Server Updates"
		},
		{
			action: "dequeue",
			disabled: (data: Array<AugmentedServer>): boolean => !data.every(serverIsCache),
			multiRow: true,
			name: "Clear Queued Updates"
		}
	];

	/** A subject that child components can subscribe to for access to the fuzzy search query text */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl: FormControl = new FormControl("");

	constructor(private readonly api: ServerService,
		public readonly auth: CurrentUserService,
		private readonly route: ActivatedRoute,
		private readonly navSvc: NavigationService,
		private readonly cdn: CDNService,
		private readonly dialog: MatDialog) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.navSvc.headerTitle.next("Servers");
		this.cdns = this.cdn.getCDNs();
	}

	/** Initializes table data, loading it from Traffic Ops. */
	public ngOnInit(): void {
		this.reloadServers();

		this.route.queryParamMap.subscribe(
			m => {
				const search = m.get("search");
				if (search) {
					this.fuzzControl.setValue(decodeURIComponent(search));
					this.fuzzySubject.next(search);
					this.fuzzySubject.next(this.fuzzControl.value);
				}
			}
		);
	}

	/** Update the URL's 'search' query parameter for the user's search input. */
	public updateURL(): void {
		this.fuzzySubject.next(this.fuzzControl.value);
	}

	/**
	 * Handles user selection of a more menu action button.
	 *
	 * @param action The emitted more menu button action event.
	 */
	public async handleMoreMenu(action: string): Promise<void> {
		const data: CollectionChoiceDialogData<number> = {
			collection: (await this.cdns).map(cdn => ({label: cdn.name, value: cdn.id})),
			label: "Please Select a CDN",
			message: "",
			title: "Queue Server Updates"
		};
		switch(action) {
			case "dequeue":
				data.title = "Clear Server Updates";
				break;
			case "queue":
				const ref = this.dialog.open<CollectionChoiceDialogComponent, CollectionChoiceDialogData<number>, number | false>(
					CollectionChoiceDialogComponent,
					{data}
				);
				const result = await ref.afterClosed().toPromise();
				if (typeof(result) === "number") {
					if (data.title.indexOf("Clear") > -1) {
						await this.cdn.dequeueServerUpdates(result);
					} else {
						await this.cdn.queueServerUpdates(result);
					}
				}
				break;
		}
	}

	/**
	 * Handles user selection of a context menu action item.
	 *
	 * @param action The emitted context menu action event.
	 */
	public async handleContextMenu(action: ContextMenuActionEvent<AugmentedServer>): Promise<void> {
		switch (action.action) {
			case "updateStatus":
				const dialogRef = this.dialog.open(UpdateStatusComponent, {
					data: action.data instanceof Array ? action.data : [action.data]
				});
				dialogRef.afterClosed().subscribe(result => {
					if(result) {
						this.reloadServers();
					}
				});
				break;
			case "queue":
				const queueServers = action.data instanceof Array ? action.data : [action.data];
				await Promise.all(queueServers.map(async s => this.api.queueUpdates(s)));
				await this.reloadServers();
				break;
			case "dequeue":
				const dequeueServers = action.data instanceof Array ? action.data : [action.data];
				await Promise.all(dequeueServers.map(async s => this.api.clearUpdates(s)));
				await this.reloadServers();
				break;
			default:
				throw new Error(`unknown context menu item clicked: ${action.action}`);
		}
	}

	/** Reloads the servers table data. */
	public async reloadServers(): Promise<void> {
		this.servers = this.api.getServers().then(ss => ss.map(augment));
	}
}
