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

import { CdkDragDrop, moveItemInArray } from "@angular/cdk/drag-drop";
import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import type {
	Interface,
	ResponseCacheGroup,
	ResponseCDN,
	ResponsePhysicalLocation,
	ResponseProfile,
	ResponseServer,
	ResponseStatus,
	TypeFromResponse
} from "trafficops-types";

import { CacheGroupService, CDNService, PhysicalLocationService, ProfileService, TypeService } from "src/app/api";
import { ServerService } from "src/app/api/server.service";
import { UpdateStatusComponent } from "src/app/core/servers/update-status/update-status.component";
import {
	DecisionDialogComponent,
	DecisionDialogData
} from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { IP, IP_WITH_CIDR, AutocompleteValue } from "src/app/utils";

/**
 * ServerDetailsComponent is the controller for a server's "details" page.
 */
@Component({
	selector: "tp-server-details",
	styleUrls: ["./server-details.component.scss"],
	templateUrl: "./server-details.component.html",
})
export class ServerDetailsComponent implements OnInit {

	/**
	 * Tracks whether the form is for creating a new server ('true') or editing an existing one ('false').
	 */
	public isNew = false;
	/**
	 * The server being edited/created.
	 */
	public server!: ResponseServer;
	/**
	 * A Regular Expression that matches valid IP addresses - and allows IPv4 addresses to have CIDR-notation network prefixes.
	 */
	public validIPPattern = IP_WITH_CIDR;
	/**
	 * A Regular Expression that matches valid IP addresses.
	 */
	public validGatewayPattern = IP;
	/**
	 * Icon for the "change status" button.
	 *
	 * @returns Material icon name
	 */
	public statusChangeIcon(): string {
		if (this.isNew || !this.server.status) {
			return "toggle_on";
		}
		if (this.server.status === "ONLINE" || this.server.status === "REPORTED") {
			return "toggle_on";
		}
		return "toggle_off";
	}

	/**
	 * The set of all Cache Groups.
	 */
	public cacheGroups = new Array<ResponseCacheGroup>();
	/**
	 * The set of all CDNs.
	 */
	public cdns = new Array<ResponseCDN>();
	/**
	 * The set of all Physical Locations.
	 */
	public physicalLocations = new Array<ResponsePhysicalLocation>();
	/**
	 * The set of all Profiles.
	 */
	public profiles = new Array<ResponseProfile>();
	/**
	 * The set of all Statuses.
	 */
	public statuses = new Array<ResponseStatus>();
	/**
	 * The set of all Types that can be applied to a server.
	 */
	public types = new Array<TypeFromResponse>();

	public autocompleteNew = AutocompleteValue.NEW_PASSWORD;

	/**
	 * Constructor.
	 */
	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly serverService: ServerService,
		private readonly cacheGroupService: CacheGroupService,
		private readonly cdnService: CDNService,
		private readonly profileService: ProfileService,
		private readonly typeService: TypeService,
		private readonly physlocService: PhysicalLocationService,
		private readonly navSvc: NavigationService,
		private readonly dialog: MatDialog,
		private readonly log: LoggingService,
	) {
	}

	/**
	 * Initializes the controller based on route query parameters.
	 */
	public ngOnInit(): void {
		const handleErr = (obj: string): (e: unknown) => void =>
			(e: unknown): void => {
				this.log.error(`Failed to get ${obj}:`, e);
			};

		this.cacheGroupService.getCacheGroups().then(
			cgs => {
				this.cacheGroups = cgs;
			}
		);
		this.cdnService.getCDNs().then(
			cdns => {
				this.cdns = Array.from(cdns.values());
			}
		);
		this.serverService.getStatuses().then(
			statuses => {
				this.statuses = statuses;
			}
		).catch(handleErr("Statuses"));
		this.profileService.getProfiles().then(
			profiles => {
				this.profiles = profiles;
			}
		).catch(handleErr("Profiles"));
		this.typeService.getServerTypes().then(
			types => {
				this.types = types;
			}
		);
		this.physlocService.getPhysicalLocations().then(
			physlocs => {
				this.physicalLocations = physlocs;
			}
		).catch(handleErr("Physical Locations"));

		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			this.log.error("missing required route parameter 'id'");
			return;
		}

		this.isNew = ID === "new";

		if (!this.isNew) {
			this.serverService.getServers(Number(ID)).then(
				s => {
					this.server = s;
					this.updateTitleBar();
				}
			).catch(
				e => {
					this.log.error(`Failed to get server #${ID}:`, e);
				}
			);
		} else {
			this.server = {
				cachegroup: "",
				cachegroupId: 0,
				cdnId: 0,
				cdnName: "",
				domainName: "",
				guid: null,
				hostName: "",
				httpsPort: null,
				id: 0,
				iloIpAddress: null,
				iloIpGateway: null,
				iloIpNetmask: null,
				iloPassword: null,
				iloUsername: null,
				interfaces: [{
					ipAddresses: [{
						address: "",
						gateway: null,
						serviceAddress: true
					}],
					maxBandwidth: null,
					monitor: true,
					mtu: null,
					name: "",
				}],
				lastUpdated: new Date(),
				mgmtIpAddress: null,
				mgmtIpGateway: null,
				mgmtIpNetmask: null,
				offlineReason: null,
				physLocation: "",
				physLocationId: 0,
				profileNames: [],
				rack: null,
				revalPending: false,
				routerHostName: null,
				routerPortName: null,
				status: "",
				statusId: 0,
				statusLastUpdated: null,
				tcpPort: null,
				type: "",
				typeId: 0,
				updPending: false,
				xmppId: ""
			};
			this.updateTitleBar();
		}
	}

	/**
	 * Updates the headerTitle based on current server state.
	 *
	 * @private
	 */
	private updateTitleBar(): void {
		if (this.isNew) {
			this.navSvc.headerTitle.next("New Server");
		} else {
			this.navSvc.headerTitle.next(`Server: ${this.server.hostName}`);
		}
	}

	/**
	 * Handles form submittal, either creating or updating the server as appropriate.
	 *
	 * @param e The raw submittal event; its default is prevented and its propagation is halted.
	 */
	public submit(e: Event): void {
		e.preventDefault();
		e.stopPropagation();
		if (this.isNew) {
			this.serverService.createServer(this.server).then(
				s => {
					if (!s.id) {
						throw new Error("Traffic Ops returned server with no ID");
					}
					this.isNew = false;
					this.server = s;
					this.router.navigate(["server", s.id]);
					this.updateTitleBar();
				},
				err => {
					this.log.error("failed to create server:", err);
				}
			);
		} else {
			this.serverService.updateServer(this.server).then(
				responseServer => {
					this.server = responseServer;
					this.updateTitleBar();
				},
				err => {
					this.log.error(`failed to update server: ${err}`);
				}
			);
		}
	}
	/**
	 * Deletes the Server.
	 */
	public delete(): void {
		if (this.isNew) {
			this.log.error("Unable to delete new Cache Group");
			return;
		}
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(
			DecisionDialogComponent,
			{
				data: {
					message: `Are you sure you want to delete Server ${this.server.hostName} (#${this.server.id})?`,
					title: "Confirm Delete"
				}
			}
		);
		ref.afterClosed().subscribe(result => {
			if (result) {
				this.serverService.deleteServer(this.server);
				this.router.navigate(["core/servers"]);
			}
		});
	}

	/**
	 * Handles when a profile list item is 'dropped'
	 *
	 * @param $event The Drop event that is emitted.
	 */
	public drop($event: CdkDragDrop<string[]>): void {
		moveItemInArray(this.server.profileNames, $event.previousIndex, $event.currentIndex);
	}

	/**
	 * Queues updates for the server
	 */
	public async queue(): Promise<void> {
		this.serverService.queueUpdates(this.server).then(result => {
			if(result.action === "queue") {
				this.server.updPending = true;
			}
		},
		err => {
			this.log.error(`failed to queue updates: ${err}`);
		});
	}

	/**
	 * Dequeues updates for the server
	 */
	public async dequeue(): Promise<void> {
		this.serverService.clearUpdates(this.server).then(result => {
			if(result.action === "dequeue") {
				this.server.updPending = false;
			}
		},
		err => {
			this.log.error(`failed to dequeue updates: ${err}`);
		});
	}

	/**
	 * Adds a new network interface to the server.
	 *
	 * @param e The triggering DOM event; its propagation is stopped.
	 */
	public addInterface(e: MouseEvent): void {
		e.stopPropagation();
		const newInf = {
			ipAddresses: [],
			maxBandwidth: null,
			monitor: false,
			mtu: 1500,
			name: ""
		};

		this.server.interfaces.push(newInf);
	}

	/**
	 * Returns a user-friendly name for an interface.
	 *
	 * @param inf The Interface to get the name from
	 * @returns Friendly interface name
	 */
	public getInterfaceName(inf: Interface): string {
		return inf.name === "" ? "<un-named>" : inf.name;
	}

	/**
	 * Finds the ID of a given profile name
	 *
	 * @param profileName The profileName to find the id of.
	 * @returns Profile id
	 */
	public profileNameToId(profileName: string): number {
		return (this.profiles.find(p => p.name === profileName) ?? {id: -1}).id;
	}

	/**
	 * Adds a new IP address to the server.
	 *
	 * @param event The triggering DOM event; its propagation is stopped.
	 * @param inf The specific network interface to which to add the new IP address.
	 */
	public addIP(event: MouseEvent, inf: Interface): void {
		event.stopPropagation();
		inf.ipAddresses.push({
			address: "",
			gateway: null,
			serviceAddress: false
		});
	}

	/**
	 * Removes an IP address from the server.
	 *
	 * @param event The triggering DOM event; its propagation is stopped.
	 * @param inf The specific network interface from which to remove an IP address.
	 * @param ip The index in the `ipAddresses` of `inf` to delete.
	 */
	public deleteIP(event: MouseEvent, inf: Interface, ip: number): void {
		event.stopPropagation();
		inf.ipAddresses.splice(ip, 1);
	}

	/**
	 * Removes a network interface from the server.
	 *
	 * @param e The triggering DOM event; its propagation is stopped.
	 * @param inf The index of the interface to remove.
	 */
	public deleteInterface(e: MouseEvent, inf: number): void {
		e.stopPropagation();
		this.server.interfaces.splice(inf, 1);
	}

	/**
	 * Tells whether the edited server is a cache server.
	 *
	 * @returns 'true' if the edited/new server is a cache server, false otherwise.
	 */
	public isCache(): boolean {
		if (!this.server.type) {
			return false;
		}
		return this.server.type.startsWith("EDGE") || this.server.type.startsWith("MID");
	}

	/**
	 * Changes the server's status.
	 *
	 * @param e The click event that triggered this handler.
	 * @throws {Error} when trying to update server status if the server doesn't exist yet (`this.isNew === true`).
	 */
	public changeStatus(e: MouseEvent): void {
		e.stopPropagation();
		if (this.isNew) {
			throw new Error("cannot update the status of a server that doesn't exist yet");
		}
		const ref = this.dialog.open(UpdateStatusComponent, {
			data: [this.server]
		});
		ref.afterClosed().subscribe(res => {
			if (res) {
				this.serverService.getServers(this.server.id).then(
					s => this.server = s
				).catch(
					err => {
						this.log.error("Failed to reload servers:", err);
					}
				);
			}
		});
	}
}
