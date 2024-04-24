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

import { Component } from "@angular/core";
import { FormControl, FormGroup } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import { RequestStatus, ResponseStatus } from "trafficops-types";

import { ServerService } from "src/app/api";
import {
	DecisionDialogComponent,
	DecisionDialogData
} from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
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
	public loading = true;

	/** All details of status requested */
	public statusDetails!: ResponseStatus;

	/** Reactive form intialized to creat / edit status details */
	public statusDetailsForm = new FormGroup({
		description: new FormControl("", {nonNullable: true}),
		name: new FormControl("", {nonNullable: true}),
	});

	constructor(
		private readonly api: ServerService,
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
	) {
		// Getting id from the route
		const id = this.route.snapshot.paramMap.get("id");

		/**
		 * Initializes table data, loading it from Traffic Ops.
		 * we check whether params is a number if not we shall assume user wants to add a new status.
		 */
		if (id && id !== "new") {
			this.getStatusDetails(id);
		} else {
			this.navSvc.headerTitle.next("New Status");
			this.new = true;
			this.loading = false;
		}
	}

	/**
	 * Get status details for the id
	 * patch the form with status details
	 *
	 * @param id ID of the status
	 */
	public async getStatusDetails(id: string | number): Promise<void> {
		this.statusDetails = await this.api.getStatuses(Number(id));

		// Set page title with status Name
		this.setTitle();

		// Patch the form with existing data we got from service requested above.
		this.statusDetailsForm.setValue({
			description: this.statusDetails.description ? this.statusDetails.description : "",
			name: this.statusDetails.name
		});

		this.loading = false;
	}

	/**
	 * Sets the headerTitle based on current Status state.
	 *
	 * @private
	 */
	private setTitle(): void {
		const title = this.new ? "New Status" : `Status: ${this.statusDetails.name}`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * On submitting the form we check for whether we are performing Create or Edit
	 *
	 * @param event HTML form submission event.
	 */
	public async onSubmit(event: Event): Promise<void> {
		event.preventDefault();
		event.stopPropagation();

		if (this.statusDetailsForm.valid) {
			if (this.new) {
				const newData: RequestStatus = {
					description: this.statusDetailsForm.controls.description.value,
					name: this.statusDetailsForm.controls.name.value
				};
				this.statusDetails = await this.api.createStatus(newData);
				this.new = false;
				await this.router.navigate(["core/statuses", this.statusDetails.id]);
			} else {
				const editData: ResponseStatus = {
					description: this.statusDetailsForm.controls.description.value,
					id: this.statusDetails.id,
					lastUpdated: this.statusDetails.lastUpdated,
					name: this.statusDetailsForm.controls.name.value
				};
				this.statusDetails = await this.api.updateStatusDetail(editData);
			}
			this.setTitle();
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
				this.router.navigate(["core/statuses"]);
			}
		});
	}
}
