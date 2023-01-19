/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import { Location } from "@angular/common";
import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { ResponsePhysicalLocation, ResponseRegion } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * PhysLocDetailComponent is the controller for the physloc add/edit form
 */
@Component({
	selector: "tp-phys-loc-detail",
	styleUrls: ["./phys-loc-detail.component.scss"],
	templateUrl: "./phys-loc-detail.component.html"
})
export class PhysLocDetailComponent implements OnInit {
	public new = false;

	public physLocation!: ResponsePhysicalLocation;
	public regions!: Array<ResponseRegion>;

	constructor(private readonly route: ActivatedRoute, private readonly cacheGroupService: CacheGroupService,
		private readonly location: Location, private readonly dialog: MatDialog, private readonly navSvc: NavigationService) {
	}

	/**
	 * Angular lifecycle hook.
	 */
	public async ngOnInit(): Promise<void> {
		this.regions = await this.cacheGroupService.getRegions();
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			console.error("missing required route parameter 'id'");
			return;
		}

		if (ID === "new") {
			this.navSvc.headerTitle.next("New Physical Location");
			this.new = true;
			this.physLocation = {
				address: "",
				city: "",
				comments: null,
				email: null,
				id: -1,
				lastUpdated: new Date(),
				name: "",
				phone: null,
				poc: null,
				region: null,
				regionId: -1,
				shortName: "",
				state: "",
				zip: ""
			};
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			console.error("route parameter 'id' was non-number:", ID);
			return;
		}

		this.physLocation = await this.cacheGroupService.getPhysicalLocations(numID);
		this.navSvc.headerTitle.next(`Physical Location: ${this.physLocation.name}`);
	}

	/**
	 * Deletes the current physLocation.
	 */
	public async deletePhysicalLocation(): Promise<void> {
		if (this.new) {
			console.error("Unable to delete new physLocation");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete physical location ${this.physLocation.name} with id ${this.physLocation.id}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.cacheGroupService.deletePhysicalLocation(this.physLocation.id);
				this.location.back();
			}
		});
	}

	/**
	 * Submits new/updated physLocation.
	 *
	 * @param e HTML click event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if(this.new) {
			this.physLocation = await this.cacheGroupService.createPhysicalLocation(this.physLocation);
			this.new = false;
		} else {
			this.physLocation = await this.cacheGroupService.updatePhysicalLocation(this.physLocation);
		}
	}
}
