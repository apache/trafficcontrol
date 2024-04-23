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
import { ResponseCoordinate } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * CoordinateDetailsComponent is the controller for the coordinate add/edit form.
 */
@Component({
	selector: "tp-coordinates-detail",
	styleUrls: ["./coordinate-detail.component.scss"],
	templateUrl: "./coordinate-detail.component.html"
})
export class CoordinateDetailComponent implements OnInit {
	public new = false;
	public coordinate!: ResponseCoordinate;

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
			this.coordinate = {
				id: -1,
				lastUpdated: new Date(),
				latitude: 0,
				longitude: 0,
				name: ""
			};
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			this.log.error("route parameter 'id' was non-number:", ID);
			return;
		}

		this.coordinate = await this.cacheGroupService.getCoordinates(numID);
		this.setTitle();
	}

	/**
	 * Sets the headerTitle based on current Coordinate state.
	 *
	 * @private
	 */
	private setTitle(): void {
		const title = this.new ? "New Coordinate" : `Coordinate: ${this.coordinate.name}`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * Deletes the current coordinate.
	 */
	public async deleteCoordinate(): Promise<void> {
		if (this.new) {
			this.log.error("Unable to delete new coordinate");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete coordinate ${this.coordinate.name} with id ${this.coordinate.id}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.cacheGroupService.deleteCoordinate(this.coordinate.id);
				this.router.navigate(["core/coordinates"]);
			}
		});
	}

	/**
	 * Submits new/updated coordinate.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if(this.new) {
			this.coordinate = await this.cacheGroupService.createCoordinate(this.coordinate);
			this.new = false;
			await this.router.navigate(["core/coordinates", this.coordinate.id]);
		} else {
			this.coordinate = await this.cacheGroupService.updateCoordinate(this.coordinate);
		}
		this.setTitle();
	}

}
