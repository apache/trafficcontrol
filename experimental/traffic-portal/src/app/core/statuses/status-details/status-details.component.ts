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
import { Component, OnInit } from "@angular/core";
import { FormControl, FormGroup, Validators } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
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
export class StatusDetailsComponent implements OnInit {

	/** Status ID expected from the route param using which we identify whether we are creating new status or load existing status */
	public id: string | number | null = null;

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
		private readonly router: Router,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
	) {
		// Getting id from the route
		this.id = this.route.snapshot.paramMap.get("id");

		// Show error if route param is null 
		if (this.id === null) {
			console.error("missing required route parameter 'id'");
			return;
		}

		/** 
		 * Initializes table data, loading it from Traffic Ops.
		 * we check whether params is a number if not we shall assume user wants to add a new status.
		 */
		if (this.id !== "new") {
			this.loading = true;
			this.statusDetailsForm.addControl("id", new FormControl(""));
			this.statusDetailsForm.addControl("lastUpdated", new FormControl(""));
			this.getStatusDetails();
		} else {
			this.navSvc.headerTitle.next("New Status");
		}
	}

	
	public ngOnInit(): void {
	}

	/**
	 * Reloads the servers table data.
	 *
	 * @param id is the id passed in route for this page if this is a edit view.
	 */
	public async getStatusDetails(): Promise<void> {
		const id = Number(this.id); // id Type 'null' is not assignable to type 'string'
		this.statusDetails = await this.api.getStatuses(id);

		// Set page title with status Name
		this.navSvc.headerTitle.next(`Status #${this.statusDetails.name}`);

		// Patch the form with existing data we got from service requested above.
		this.statusDetailsForm.patchValue(this.statusDetails);
		this.loading = false;
	}

	/**
	 * On submitting the form we check for whether we are performing Create or Edit
	 * @param event The DOM form submission event.
	 */
	public onSubmit(event: Event): void {
		event.preventDefault();
		event.stopPropagation();

		if (this.statusDetailsForm.invalid) {
			return;
		}
		if (this.id === "new") {
			this.createStatus();

		} else {
			this.updateStatus();
		}
	}

	/**
	 * For Creating a new status
	 */
	public createStatus(): void {
		this.api.createStatus(this.statusDetailsForm.value).then((res: ResponseStatus) => {
			if (res) {
				this.statusDetails = res;
				this.router.navigate(["/core/statuses"]);
			}
		});
	}

	/**
	 * For updating the Status
	 */
	public updateStatus(): void {
		this.api.updateStatusDetail(this.statusDetailsForm.value, Number(this.id));
	}

	/**
	 * Deleteting status
	 */
	public async deleteStatus(): Promise<void> {
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(DecisionDialogComponent, {
			data: {
				message: `This action CANNOT be undone. This will permanently delete '${this.statusDetails?.name}'.`,
				title: `Delete Status: ${this.statusDetails?.name}`
			}
		});

		ref.afterClosed().subscribe(result => {
			if (result) {
				const id = Number(this.id);
				this.api.deleteStatus(id).then(() => {
					this.router.navigate(["/core/statuses"]);
				});
			}
		});
	}
}
