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
import { FormBuilder, FormControl, FormGroup, Validators } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import { ResponseStatus } from "trafficops-types";

import { StatusesService } from "src/app/api/statuses.service";
import { DecisionDialogComponent, DecisionDialogData } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";

@Component({
	selector: "tp-status-details",
	templateUrl: "./status-details.component.html",
	styleUrls: ["./status-details.component.scss"]
})
export class StatusDetailsComponent implements OnInit {

	id: string | null = null;
	statusDetails: ResponseStatus | null = null;
	statusDetailsForm!: FormGroup;
	loading = false;
	submitting = false;
	submitted = false;

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly fb: FormBuilder,
		private readonly dialog: MatDialog,
		private readonly statusesService: StatusesService) { }

	ngOnInit(): void {
		// Form is built here
		this.statusDetailsForm = this.fb.group({
			name: ["", Validators.required],
			description: ["", Validators.required],
		});

		// Getting id from the route
		this.id = this.route.snapshot.paramMap.get("id");

		// we check whether params is a number if not we shall assume user wants to add a new status.
		if (!this.isNew) {
			this.loading = true;
			this.statusDetailsForm.addControl("id", new FormControl(""));
			this.statusDetailsForm.addControl("lastUpdated", new FormControl(""));
			this.getStatusDetails();
		}
	}

	/*
   * Reloads the servers table data.
   * @param id is the id passed in route for this page if this is a edit view.
  */
	async getStatusDetails(): Promise<void> {
		const id = this.id as string; // id Type 'null' is not assignable to type 'string'
		this.statusDetails = await this.statusesService.getStatuses(id);
		const data: ResponseStatus = {
			name: this.statusDetails.name,
			description: this.statusDetails.description,
			lastUpdated: new Date(),
			id: this.statusDetails.id
		};
		this.statusDetailsForm.patchValue(data);
		this.loading = false;
	}

	// On submitting the form we check for whether we are performing Create or Edit
	onSubmit() {
		if (this.isNew) {
			this.createStatus();

		} else {
			this.updateStatus();
		}
	}

	// For Creating a new status
	createStatus() {
		this.statusesService.createStatus(this.statusDetailsForm.value).then((res: any) => {
			if (res) {
				this.id = res?.id;
				this.router.navigate([`/core/statuses/${this.id}`]);
			}
		});
	}

	// For updating the Status
	updateStatus() {
		this.statusesService.updateStatus(this.statusDetailsForm.value, Number(this.id));
	}

	// Deleteting status
	async deleteStatus() {
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(DecisionDialogComponent, {
			data: {
				message: `This action CANNOT be undone. This will permanently delete '${this.statusDetails?.name}'.`,
				title: `Delete Status: ${this.statusDetails?.name}`
			}
		});

		if (await ref.afterClosed().toPromise()) {
			const id = Number(this.id);
			this.statusesService.deleteStatus(id).then(() => {
				this.router.navigate(["/core/statuses"]);
			});
		}

	}

	// Title for the page
	get title(): string {
		return this.isNew ? "Add New Status" : "Edit Status";
	}

	// Checking for params to ensure given id is a number
	get isNew() {
		return this.id === "new" && isNaN(Number(this.id));
	}
}
