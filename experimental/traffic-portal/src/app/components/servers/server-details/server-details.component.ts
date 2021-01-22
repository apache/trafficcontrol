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
import { ActivatedRoute } from "@angular/router";
import { faClock, faMinus, faPlus, faToggleOff, faToggleOn, IconDefinition } from "@fortawesome/free-solid-svg-icons";
import { faClock as hollowClock } from "@fortawesome/free-regular-svg-icons";
import { CacheGroup, CDN, DUMMY_SERVER, Interface, Profile, Server, Status } from "src/app/models";
import { CacheGroupService, CDNService, ProfileService, ServerService } from "src/app/services/api";
import { IP, IP_WITH_CIDR } from "src/app/utils";

@Component({
	selector: "tp-server-details",
	styleUrls: ["./server-details.component.scss"],
	templateUrl: "./server-details.component.html",
})
export class ServerDetailsComponent implements OnInit {

	/**
	 *
	 */
	public isNew = false;
	/**
	 *
	 */
	public server: Server;
	/**
	 *
	 */
	public validIPPattern = IP_WITH_CIDR;
	/**
	 *
	 */
	public validGatewayPattern = IP;

	/**
	 *
	 */
	public get title(): string {
		if (this.isNew) {
			return "New Server";
		}
		return `Server #${this.server.id}`;
	}

	/**
	 *
	 */
	public hideILO = false;
	/**
	 *
	 */
	public hideManagement = false;
	/**
	 *
	 */
	public hideInterfaces = false;

	/**
	 *
	 */
	public addIcon = faPlus;
	/**
	 *
	 */
	public removeIcon = faMinus;
	/**
	 *
	 */
	public clearUpdatesIcon = faClock;
	/**
	 *
	 */
	public updateIcon = hollowClock;
	/**
	 *
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
	 *
	 */
	public cacheGroups = new Array<CacheGroup>();
	/**
	 *
	 */
	public cdns = new Array<CDN>();
	// public physicalLocations = new Array<PhysicalLocation>();
	/**
	 *
	 */
	public profiles = new Array<Profile>();
	/**
	 *
	 */
	public statuses = new Array<Status>();

	/**
	 *
	 */
	constructor(
		private readonly route: ActivatedRoute,
		private readonly serverService: ServerService,
		private readonly cacheGroupService: CacheGroupService,
		private readonly cdnService: CDNService,
		private readonly profileService: ProfileService
	) {
		this.server = DUMMY_SERVER;
	}

	/**
	 *
	 */
	public serverJSON(): string {
		return JSON.stringify(this.server, null, "\t");
	}

	/**
	 * Initializes the controller based on route query parameters.
	 */
	public ngOnInit(): void {
		this.cacheGroupService.getCacheGroups().subscribe(
			cgs => {
				this.cacheGroups = cgs;
			}
		);
		this.cdnService.getCDNs().subscribe(
			cdns => {
				this.cdns = Array.from(cdns.values());
			}
		);
		this.serverService.getStatuses().subscribe(
			statuses => {
				this.statuses = statuses;
			}
		);
		this.profileService.getProfiles().subscribe(
			profiles => {
				this.profiles = profiles;
			}
		);

		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			console.error("missing required route parameter 'id'");
			return;
		}

		this.isNew = ID === "new";

		if (!this.isNew) {
			this.serverService.getServers(Number(ID)).subscribe(
				s => {
					this.server = s;
				}
			);
		}
	}

	/**
	 *
	 */
	public submit(e: Event): void {
		e.preventDefault();
		e.stopPropagation();
		if (this.isNew) {
			this.serverService.createServer(this.server).subscribe(
				s => {
					this.server = s;
					this.isNew = false;
				},
				err => {
					console.error("failed to create server:", err);
				}
			);
		}
	}

	/**
	 *
	 */
	public log(): void {
		console.log(this);
	}

	/**
	 *
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
	 *
	 */
	public addIP(inf: Interface): void {
		inf.ipAddresses.push({
			address: "",
			gateway: null,
			serviceAddress: false
		});
	}

	/**
	 *
	 */
	public deleteIP(inf: Interface, ip: number): void {
		inf.ipAddresses.splice(ip, 1);
	}

	/**
	 *
	 */
	public deleteInterface(inf: number): void {
		this.server.interfaces.splice(inf, 1);
	}

	/**
	 *
	 */
	public isCache(): boolean {
		if (!this.server.type) {
			return false;
		}
		return this.server.type.startsWith("EDGE") || this.server.type.startsWith("MID");
	}

}
