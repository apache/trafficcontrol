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
import { Component } from "@angular/core";
import { FormControl, FormGroup, Validators } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { ResponseStatus } from "trafficops-types";

import { ServerService } from "src/app/api";
import { DecisionDialogComponent, DecisionDialogData } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * StatusDetailsComponent is the controller for a status "details" page.
 */
@Component({
	selector: "tp-status-details",
	styleUrls: ["./status-details.component.scss"],
	templateUrl: "./status-details.component.html",
})
export class StatusDetailsComponent {
	public new = false;

	/** Loader status for the actions */
	public loading = false;

	/** All details of status requested */
	public statusDetails!: ResponseStatus;

	/** Reactive form intialized to creat / edit status details */
	public statusDetailsForm: FormGroup = new FormGroup({
		description: new FormControl("", Validators.required),
		name: new FormControl("", Validators.required),
	});

	/**
	 * Constructor.
	 *
	 * @param api The Servers API which is used to provide row data.
	 * @param route A reference to the route of this view which is used to get the 'id' query parameter of status.
	 * @param router Angular router
	 * @param dialog Dialog manager
	 * @param fb Form builder
	 * @param navSvc Manages the header
	 */
	constructor(
		private readonly api: ServerService,
		private readonly route: ActivatedRoute,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService, private readonly location: Location,
	) {
		// Getting id from the route
		const id = this.route.snapshot.paramMap.get("id");

		/**
		 * Initializes table data, loading it from Traffic Ops.
		 * we check whether params is a number if not we shall assume user wants to add a new status.
		 */
		if (id !== "new") {
			this.loading = true;
			this.statusDetailsForm.addControl("id", new FormControl(""));
			this.statusDetailsForm.addControl("lastUpdated", new FormControl(""));
			this.getStatusDetails();
		} else {
			this.navSvc.headerTitle.next("New Status");
			this.new = true;
		}
	}

	/**
	 * Get status details for the id
	 * patch the form with status details
	 */
	public async getStatusDetails(): Promise<void> {
		this.statusDetails = await this.api.getStatuses(this.statusDetails.id);

		// Set page title with status Name
		this.navSvc.headerTitle.next(`Status #${this.statusDetails.name}`);

		// Patch the form with existing data we got from service requested above.
		this.statusDetailsForm.patchValue(this.statusDetails);
		this.loading = false;
	}

	/**
	 * On submitting the form we check for whether we are performing Create or Edit
	 *
	 * @param event HTML form submission event.
	 */
	public async onSubmit(event: Event): Promise<void>  {
		event.preventDefault();
		event.stopPropagation();

		if (this.statusDetailsForm.valid) {
			if (this.new) {
				this.statusDetails = await this.api.createStatus(this.statusDetailsForm.value);
				this.location.back();
			} else {
				this.statusDetails = await this.api.updateStatusDetail(this.statusDetailsForm.value);
			}
		}
	}

	/**
	 * Deleteting status
	 */
	public async deleteStatus(): Promise<void> {
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(DecisionDialogComponent, {
			data: {
				message: `This action CANNOT be undone. This will permanently delete '${this.statusDetails.name}'.`,
				title: `Delete Status: ${this.statusDetails.name}`
			}
		});

		ref.afterClosed().subscribe(result => {
			if (result) {
				this.api.deleteStatus(this.statusDetails.id);
				this.location.back();
			}
		});
	}
}
