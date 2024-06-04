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
		private readonly router: Router,
		private readonly cacheGroupService: CacheGroupService,
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

		this.new = ID === "new";

		if (this.new) {
			this.setTitle();
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
		this.setTitle();
	}

	/**
	 * Sets the headerTitle based on current Division state.
	 *
	 * @private
	 */
	private setTitle(): void {
		const title = this.new ? "New Division" : `Division: ${this.division.name}`;
		this.navSvc.headerTitle.next(title);
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
				this.router.navigate(["core/divisions"]);
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
			await this.router.navigate(["core/divisions", this.division.id]);
		} else {
			this.division = await this.cacheGroupService.updateDivision(this.division);
		}
		this.setTitle();
	}

}
