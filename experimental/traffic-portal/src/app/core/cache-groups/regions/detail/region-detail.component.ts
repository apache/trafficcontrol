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
import { Location } from "@angular/common";
import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { ResponseDivision, ResponseRegion } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * RegionDetailsComponent is the controller for the region add/edit form.
 */
@Component({
	selector: "tp-regions-detail",
	styleUrls: ["./region-detail.component.scss"],
	templateUrl: "./region-detail.component.html"
})
export class RegionDetailComponent implements OnInit {
	public new = false;
	public region!: ResponseRegion;
	public divisions!: Array<ResponseDivision>;

	constructor(private readonly route: ActivatedRoute, private readonly cacheGroupService: CacheGroupService,
		private readonly location: Location, private readonly dialog: MatDialog,
		private readonly header: NavigationService) {
	}

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		this.divisions = await this.cacheGroupService.getDivisions();
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			console.error("missing required route parameter 'id'");
			return;
		}

		if (ID === "new") {
			this.header.headerTitle.next("New Region");
			this.new = true;
			this.region = {
				division: -1,
				divisionName: "",
				id: -1,
				lastUpdated: new Date(),
				name: ""
			};
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			console.error("route parameter 'id' was non-number:", ID);
			return;
		}

		this.region = await this.cacheGroupService.getRegions(numID);
		this.header.headerTitle.next(`Region: ${this.region.name}`);
	}

	/**
	 * Deletes the current region.
	 */
	public async deleteRegion(): Promise<void> {
		if (this.new) {
			console.error("Unable to delete new region");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete region ${this.region.name} with id ${this.region.id}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.cacheGroupService.deleteRegion(this.region.id);
				this.location.back();
			}
		});
	}

	/**
	 * Submits new/updated region.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if(this.new) {
			this.region = await this.cacheGroupService.createRegion(this.region);
			this.new = false;
		} else {
			this.region = await this.cacheGroupService.updateRegion(this.region);
		}
	}

}
