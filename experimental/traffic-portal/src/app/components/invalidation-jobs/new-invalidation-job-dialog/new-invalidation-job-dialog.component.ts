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
import { Component, Inject } from "@angular/core";
import { FormControl } from "@angular/forms";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { Subject } from "rxjs";

import { InvalidationJobService } from "src/app/services/api";

/**
 * This is the controller for the dialog box that opens when the user creates
 * a new Content Invalidation Job.
 */
@Component({
	selector: "tp-new-invalidation-job-dialog",
	styleUrls: ["./new-invalidation-job-dialog.component.scss"],
	templateUrl: "./new-invalidation-job-dialog.component.html",
})
export class NewInvalidationJobDialogComponent {
	/** The minimum Start Date that may be selected. */
	public readonly startMin = new Date();
	/** The minimum Start Time that may be selected. */
	public startMinTime: string;

	/** The date - but not time! - at which the new job will come into effect. */
	public startDate = new Date();
	/** Control for users to enter new content invalidation jobs. */
	public regexp = new FormControl("/");
	/** Control for users to enter a new job's TTL. */
	public ttl = new FormControl(178);
	/** Control for users to enter the starting time for a new job. */
	public startTime = new FormControl("");

	/** A subscribable that tracks whether the new job's regexp is valid. */
	public readonly regexpIsValid = new Subject<string>();

	constructor(
		private readonly dialogRef: MatDialogRef<NewInvalidationJobDialogComponent>,
		private readonly jobAPI: InvalidationJobService,
		@Inject(MAT_DIALOG_DATA) private readonly dsID: number
	) {
		this.startDate.setDate(this.startDate.getDate()+1);
		const hours = String(this.startMin.getHours()).padStart(2, "0");
		const minutes = String(this.startMin.getMinutes()).padStart(2, "0");
		this.startMinTime = `${hours}:${minutes}`;
		this.startTime.setValue(this.startMinTime);
	}

	/**
	 * Updates the minimum start time when the date changes (if the date is
	 * today the current time is the minimum time, otherwise it's 00:00).
	 */
	public dateChange(): void {
		if (
			this.startDate.getFullYear() <= this.startMin.getFullYear() &&
			this.startDate.getMonth() <= this.startMin.getMonth() &&
			this.startDate.getDate() <= this.startMin.getDate()
		) {
			const hours = String(this.startMin.getHours()).padStart(2, "0");
			const minutes = String(this.startMin.getMinutes()).padStart(2, "0");
			this.startMinTime = `${hours}:${minutes}`;
		} else {
			this.startMinTime = "00:00";
		}
	}

	/**
	 * Handles submission of the content invalidation job creation form.
	 *
	 * @param event The form submission event, which must be .preventDefault'd.
	 */
	public onSubmit(event: Event): void {
		event.preventDefault();
		event.stopPropagation();

		let re: RegExp;
		try {
			re = new RegExp(this.regexp.value);
		} catch (err) {
			this.regexpIsValid.next(`Must be a valid regular expression! (${err})`);
			return;
		}

		const startTime = new Date(this.startDate);
		const [hours, minutes] = (this.startTime.value as string).split(":").map(x=>Number(x));
		startTime.setHours(hours);
		startTime.setMinutes(minutes);

		const job = {
			deliveryService: this.dsID,
			regex: re.toString().replace(/^\/|\/$/g, "").replace("\\/", "/"),
			startTime,
			ttl: this.ttl.value
		};

		this.jobAPI.createInvalidationJob(job).then(
			r => {
				if (r) {
					this.dialogRef.close(true);
				} else {
					console.warn("failure");
				}
			},
			err => {
				console.error("error: ", err);
			}
		);
	}

	/**
	 * Closes the dialog, indicating that no new Job was created.
	 */
	public cancel(): void {
		this.dialogRef.close();
	}
}
