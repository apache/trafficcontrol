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
import { FormControl, UntypedFormControl } from "@angular/forms";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { Subject } from "rxjs";
import { JobType, ResponseInvalidationJob } from "trafficops-types";

import { InvalidationJobService } from "src/app/api";
import { LoggingService } from "src/app/shared/logging.service";

/**
 * Gets the time part of a Date as a string.
 *
 * @param d The date to convert.
 * @returns A string that represents the *local* time of the given date in
 * `HH:mm` format.
 */
export function timeStringFromDate(d: Date): string {
	const hours = String(d.getHours()).padStart(2, "0");
	const minutes = String(d.getMinutes()).padStart(2, "0");
	return `${hours}:${minutes}`;
}

/** The type of parameters passable to the dialog. */
export interface NewInvalidationJobDialogData {
	/** The ID of the Delivery Service to which the created/edited Job belongs. */
	dsID: string;
	/** If passed, the dialog will edit this Job instead of creating a new one. */
	job?: ResponseInvalidationJob;
}

/**
 * Gets the string representation of a user-entered regular expression (for
 * Content Invalidation Jobs).
 *
 * Users have a tendency to assign undue importance to '/' because of its
 * ubiquitous use in the `sed` command line utility examples and snippets
 * online. This will un-escape any '/'s that the user escaped.
 *
 * @example
 * const r = new RegExp("/.+\\/mypath\\/.+\.jpg/");
 * console.log(sanitizedRegExpString(r));
 * // Output: ".+/mypath/.+\.jpg"
 *
 * @param r A regular expression entered by a user.
 * @returns The string representation of the regexp, with unnecessary bits
 * removed.
 */
export function sanitizedRegExpString(r: RegExp): string {
	return r.toString().replace(/^\/|\/$/g, "").replace(/\\\//g, "/");
}

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
	public startMin = new Date();
	/** The minimum Start Time that may be selected. */
	public startMinTime: string;

	/** The date - but not time! - at which the new job will come into effect. */
	public startDate = new Date();
	/** Control for users to enter new content invalidation jobs. */
	public regexp = new UntypedFormControl("/");
	/** Control for users to enter a new job's TTL. */
	public ttl = new FormControl(178, {nonNullable: true});
	/** Control for users to enter the starting time for a new job. */
	public startTime = new UntypedFormControl("");

	/** A subscribable that tracks whether the new job's regexp is valid. */
	public readonly regexpIsValid = new Subject<string>();

	private readonly job: ResponseInvalidationJob | undefined;
	private readonly dsID: string;

	constructor(
		private readonly dialogRef: MatDialogRef<NewInvalidationJobDialogComponent>,
		private readonly jobAPI: InvalidationJobService,
		@Inject(MAT_DIALOG_DATA) data: NewInvalidationJobDialogData,
		private readonly log: LoggingService,
	) {
		this.job = data.job;
		if (this.job) {
			this.startDate = this.job.startTime;
			const startTime  = timeStringFromDate(this.job.startTime);
			this.startMinTime = startTime;
			this.startTime.setValue(startTime);
			this.ttl.setValue(this.job.ttlHours);
			const regexp = this.job.assetUrl.split("/", 4).slice(3).join("/") || "/";
			this.regexp.setValue(regexp);
		} else {
			this.startDate.setDate(this.startDate.getDate()+1);
			this.startMinTime = timeStringFromDate(this.startMin);
			this.startTime.setValue(this.startMinTime);
		}

		this.dsID = data.dsID;
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
			this.startMinTime = timeStringFromDate(this.startMin);
		} else {
			this.startMinTime = "00:00";
		}
	}

	/**
	 * Updates a content invalidation job based on passed pre-parsed data in
	 * combination with the component's state.
	 *
	 * @param j The Job being edited (in its original form).
	 * @param re The Job's new Regular Expression (pre-parsed from a form
	 * control).
	 * @param startTime The Job's new Start Time (pre-parsed from Form Controls).
	 */
	private async editJob(j: ResponseInvalidationJob, re: RegExp, startTime: Date): Promise<void> {
		const job = {
			...j,
			assetUrl: `${j.assetUrl.split("/").slice(0, 3).join("/")}/${sanitizedRegExpString(re)}`,
			startTime,
			ttlHours: this.ttl.value,
		};

		try {
			await this.jobAPI.updateInvalidationJob(job);
			this.dialogRef.close(true);
		} catch (e) {
			this.log.error(`failed to edit Job #${j.id}:`, e);
		}
	}

	/**
	 * Handles submission of the content invalidation job creation form.
	 *
	 * @param event The form submission event, which must be .preventDefault'd.
	 */
	public async onSubmit(event: Event): Promise<void> {
		event.preventDefault();
		event.stopPropagation();

		let re: RegExp;
		try {
			re = new RegExp(this.regexp.value);
		} catch (err) {
			return this.regexpIsValid.next(`Must be a valid regular expression! (${err instanceof Error ? err.message : err})`);
		}

		const startTime = new Date(this.startDate);
		const [hours, minutes] = (this.startTime.value as `${number}:${number}`).split(":").map(x=>Number(x));
		startTime.setHours(hours);
		startTime.setMinutes(minutes);

		if (this.job) {
			return this.editJob(this.job, re, startTime);
		}

		const job = {
			deliveryService: this.dsID,
			invalidationType: JobType.REFRESH,
			regex: re.toString().replace(/^\/|\/$/g, "").replace("\\/", "/"),
			startTime,
			ttlHours: this.ttl.value
		};

		try {
			await this.jobAPI.createInvalidationJob(job);
			this.dialogRef.close(true);
		} catch (err) {
			this.log.error("failed to create invalidation job: ", err);
		}
	}

	/**
	 * Closes the dialog, indicating that no new Job was created.
	 */
	public cancel(): void {
		this.dialogRef.close();
	}
}
