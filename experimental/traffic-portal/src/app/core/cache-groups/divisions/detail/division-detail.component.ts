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
import { ResponseDivision } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * DivisionDetailsComponent is the controller for the division add/edit form.
 */
@Component({
	selector: "tp-divisions-detail",
	templateUrl: "./division-detail.component.html"
})
export class DivisionDetailComponent implements OnInit {
	public new = false;
	public division!: ResponseDivision;

	constructor(
		private readonly route: ActivatedRoute,
		private readonly cacheGroupService: CacheGroupService,
		private readonly location: Location,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) { }

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			this.log.error("missing required route parameter 'id'");
			return;
		}

		if (ID === "new") {
			this.navSvc.headerTitle.next("New Division");
			this.new = true;
			this.division = {
				id: -1,
				lastUpdated: new Date(),
				name: ""
			};
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			this.log.error("route parameter 'id' was non-number:", ID);
			return;
		}

		this.division = await this.cacheGroupService.getDivisions(numID);
		this.navSvc.headerTitle.next(`Division: ${this.division.name}`);
	}

	/**
	 * Deletes the current division.
	 */
	public async deleteDivision(): Promise<void> {
		if (this.new) {
			this.log.error("Unable to delete new division");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete division ${this.division.name} with id ${this.division.id}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.cacheGroupService.deleteDivision(this.division.id);
				this.location.back();
			}
		});
	}

	/**
	 * Submits new/updated division.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if(this.new) {
			this.division = await this.cacheGroupService.createDivision(this.division);
			this.new = false;
		} else {
			this.division = await this.cacheGroupService.updateDivision(this.division);
		}
	}

}
