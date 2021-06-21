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

import { Component, OnInit } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { faClock, faMinus, faPlus, faToggleOff, faToggleOn, IconDefinition } from "@fortawesome/free-solid-svg-icons";
import { faClock as hollowClock } from "@fortawesome/free-regular-svg-icons";
import { CacheGroup, CDN, DUMMY_SERVER, Interface, PhysicalLocation, Profile, Server, Status, Type } from "src/app/models";
import { CacheGroupService, CDNService, ProfileService, ServerService, TypeService } from "src/app/services/api";
import { IP, IP_WITH_CIDR } from "src/app/utils";
import { PhysicalLocationService } from "src/app/services/api/PhysicalLocationService";

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
	public server: Server;
	/**
	 * A Regular Expression that matches valid IP addresses - and allows IPv4 addresses to have CIDR-notation network prefixes.
	 */
	public validIPPattern = IP_WITH_CIDR;
	/**
	 * A Regular Expression that matches valid IP addresses.
	 */
	public validGatewayPattern = IP;
	/**
	 * Controls whether or not the "change status" dialog is open
	 */
	public changeStatusDialogOpen = false;

	/**
	 * The page title.
	 */
	public get title(): string {
		if (this.isNew) {
			return "New Server";
		}
		return `Server #${this.server.id}`;
	}

	/**
	 * Tracks whether ILO details should be hidden.
	 */
	public hideILO = false;
	/**
	 * Tracks whether management interface details should be hidden.
	 */
	public hideManagement = false;
	/**
	 * Tracks whether network interface details should be hidden.
	 */
	public hideInterfaces = false;

	/**
	 * Icon for adding to a collection.
	 */
	public addIcon = faPlus;
	/**
	 * Icon for removing from a collection.
	 */
	public removeIcon = faMinus;
	/**
	 * Icon for the "clear updates" button.
	 */
	public clearUpdatesIcon = faClock;
	/**
	 * Icon for the "queue updates" button.
	 */
	public updateIcon = hollowClock;
	/**
	 * Icon for the "change status" button.
	 */
	public get statusChangeIcon(): IconDefinition {
		if (this.isNew || !this.server.status) {
			return faToggleOn;
		}
		if (this.server.status === "ONLINE" || this.server.status === "REPORTED") {
			return faToggleOn;
		}
		return faToggleOff;
	}

	/**
	 * The set of all Cache Groups.
	 */
	public cacheGroups = new Array<CacheGroup>();
	/**
	 * The set of all CDNs.
	 */
	public cdns = new Array<CDN>();
	/**
	 * The set of all Physical Locations.
	 */
	public physicalLocations = new Array<PhysicalLocation>();
	/**
	 * The set of all Profiles.
	 */
	public profiles = new Array<Profile>();
	/**
	 * The set of all Statuses.
	 */
	public statuses = new Array<Status>();
	/**
	 * The set of all Types that can be applied to a server.
	 */
	public types = new Array<Type>();

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
		private readonly physlocService: PhysicalLocationService
	) {
		this.server = DUMMY_SERVER;
	}

	/**
	 * Initializes the controller based on route query parameters.
	 */
	public ngOnInit(): void {

		const handleErr = (obj: string): (e: unknown) => void =>
			(e: unknown): void => {
				console.error(`Failed to get ${obj}:`, e);
			};
		;

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
			console.error("missing required route parameter 'id'");
			return;
		}

		this.isNew = ID === "new";

		if (!this.isNew) {
			this.serverService.getServers(Number(ID)).then(
				s => {
					this.server = s;
				}
			).catch(
				e => {
					console.error(`Failed to get server #${ID}:`, e);
				}
			);
		} else {
			this.server.interfaces = [{
				ipAddresses: [{
					address: "",
					gateway: null,
					serviceAddress: true
				}],
				maxBandwidth: null,
				monitor: false,
				mtu: null,
				name: "",
			}];
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
				},
				err => {
					console.error("failed to create server:", err);
				}
			);
		}
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
	 * Adds a new IP address to the server.
	 *
	 * @param inf The specific network interface to which to add the new IP address.
	 */
	public addIP(inf: Interface): void {
		inf.ipAddresses.push({
			address: "",
			gateway: null,
			serviceAddress: false
		});
	}

	/**
	 * Removes an IP address from the server.
	 *
	 * @param inf The specific network interface from which to remove an IP address.
	 * @param ip The index in the `ipAddresses` of `inf` to delete.
	 */
	public deleteIP(inf: Interface, ip: number): void {
		inf.ipAddresses.splice(ip, 1);
	}

	/**
	 * Removes a network interface from the server.
	 *
	 * @param inf The index of the interface to remove.
	 */
	public deleteInterface(inf: number): void {
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
		this.changeStatusDialogOpen = true;
	}

	/**
	 * Handles the completion of a server update, closing the dialog and updating the view if necessary.
	 *
	 * @param reload Whether or not the server was actually changed (and thus needs to be reloaded)
	 */
	public doneUpdatingStatus(reload: boolean): void {
		this.changeStatusDialogOpen = false;
		if (this.isNew || !this.server.id) {
			console.error("done fired on server with no ID");
			return;
		}
		if (reload) {
			this.serverService.getServers(this.server.id).then(
				s => this.server = s
			).catch(
				e => {
					console.error("Failed to reload servers:", e);
				}
			);
		}
	}

}
