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
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import type {
	RequestOrigin,
	RequestOriginResponse,
	ResponseCacheGroup,
	ResponseCoordinate,
	ResponseDeliveryService,
	ResponseProfile,
	ResponseTenant,
} from "trafficops-types";

import {
	CacheGroupService,
	DeliveryServiceService,
	OriginService,
	ProfileService,
	UserService,
} from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * OriginDetailComponent is the controller for the origin add/edit form.
 */
@Component({
	selector: "tp-origins-detail",
	styleUrls: ["./origin-detail.component.scss"],
	templateUrl: "./origin-detail.component.html",
})
export class OriginDetailComponent implements OnInit {
	public new = false;
	public origin!: RequestOriginResponse;
	public tenants = new Array<ResponseTenant>();
	public coordinates = new Array<ResponseCoordinate>();
	public cacheGroups = new Array<ResponseCacheGroup>();
	public profiles = new Array<ResponseProfile>();
	public deliveryServices = new Array<ResponseDeliveryService>();
	public protocols = new Array<string>();

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly originService: OriginService,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
		private readonly userService: UserService,
		private readonly cacheGroupService: CacheGroupService,
		private readonly profileService: ProfileService,
		private readonly dsService: DeliveryServiceService
	) {}

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		this.tenants = await this.userService.getTenants();
		this.cacheGroups = await this.cacheGroupService.getCacheGroups();
		this.coordinates = await this.cacheGroupService.getCoordinates();
		this.profiles = await this.profileService.getProfiles();
		this.deliveryServices = await this.dsService.getDeliveryServices();
		this.protocols = ["http", "https"];

		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			this.log.error("missing required route parameter 'id'");
			return;
		}

		this.new = ID === "new";

		if (this.new) {
			this.setTitle();
			this.new = true;
			this.origin = {
				cachegroup: null,
				cachegroupId: -1,
				coordinate: null,
				coordinateId: -1,
				deliveryService: null,
				deliveryServiceId: -1,
				fqdn: "",
				id: -1,
				ip6Address: null,
				ipAddress: null,
				isPrimary: null,
				lastUpdated: new Date(),
				name: "",
				port: null,
				profile: null,
				profileId: -1,
				protocol: "https",
				tenant: null,
				tenantId: -1,
			};
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			this.log.error("route parameter 'id' was non-number: ", ID);
			return;
		}
		this.origin = await this.originService.getOrigins(numID);
		this.setTitle();
	}

	/**
	 * Sets the headerTitle based on current Origin state.
	 *
	 * @private
	 */
	private setTitle(): void {
		const title = this.new ? "New Origin" : `Origin: ${this.origin.name}`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * Deletes the current origin.
	 */
	public async deleteOrigin(): Promise<void> {
		if (this.new) {
			this.log.error("Unable to delete new origin");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {
				message: `Are you sure you want to delete origin ${this.origin.name}`,
				title: "Confirm Delete",
			},
		});
		ref.afterClosed().subscribe((result) => {
			if (result) {
				this.originService.deleteOrigin(this.origin);
				this.router.navigate(["core/origins"]);
			}
		});
	}

	/**
	 * Submits new/updated origin.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if (this.new) {
			const {
				cachegroupId,
				coordinateId,
				deliveryServiceId,
				fqdn,
				ipAddress,
				ip6Address,
				name,
				port,
				protocol,
				profileId,
				tenantId,
			} = this.origin;

			const requestOrigin: RequestOrigin = {
				deliveryServiceId,
				fqdn,
				ip6Address,
				ipAddress,
				name,
				port,
				protocol,
				tenantID: tenantId,
			};

			if (coordinateId !== -1) {
				requestOrigin.coordinateId = coordinateId;
			}

			if (cachegroupId !== -1) {
				requestOrigin.cachegroupId = cachegroupId;
			}

			if (profileId !== -1) {
				requestOrigin.profileId = profileId;
			}

			this.origin = await this.originService.createOrigin(requestOrigin);
			this.new = false;
			await this.router.navigate(["core/origins", this.origin.id]);
		} else {
			this.origin = await this.originService.updateOrigin(this.origin);
		}
		this.setTitle();
	}
}
