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
import { Component, type OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { faPlus, faTrash, faPencilAlt } from "@fortawesome/free-solid-svg-icons";

import { DeliveryServiceService, InvalidationJobService } from "src/app/api";
import { defaultDeliveryService, type DeliveryService, type InvalidationJob } from "src/app/models";
import {TpHeaderService} from "src/app/shared/tp-header/tp-header.service";

import { NewInvalidationJobDialogComponent } from "./new-invalidation-job-dialog/new-invalidation-job-dialog.component";

/**
 * InvalidationJobsComponent is the controller for the page that displays the
 * content invalidation jobs running for a Delivery Service.
 */
@Component({
	selector: "invalidation-jobs",
	styleUrls: ["./invalidation-jobs.component.scss"],
	templateUrl: "./invalidation-jobs.component.html"
})
export class InvalidationJobsComponent implements OnInit {

	/** The Delivery Service for which jobs are being described. */
	public deliveryservice: DeliveryService;

	/** All of the jobs for the described Delivery Service. */
	public jobs: Array<InvalidationJob>;

	/** The current date/time when the page loads */
	public now: Date = new Date();

	/** The ID of the Delivery Service to which the displayed Jobs belong. */
	private dsID = -1;

	/** The icon for the "Create a new Job" FAB. */
	public readonly addIcon = faPlus;

	/** The icon for the Job deletion button. */
	public readonly deleteIcon = faTrash;

	/** The icon for the Job edit button. */
	public readonly editIcon = faPencilAlt;

	constructor(
		private readonly route: ActivatedRoute,
		private readonly jobAPI: InvalidationJobService,
		private readonly dsAPI: DeliveryServiceService,
		private readonly dialog: MatDialog,
		private readonly headerSvc: TpHeaderService
	) {
		this.deliveryservice = {...defaultDeliveryService};
		this.jobs = new Array<InvalidationJob>();
	}

	/**
	 * Runs initialization, fetching the jobs and Delivery Service data from
	 * Traffic Ops and setting the pageload date/time.
	 */
	public ngOnInit(): void {
		this.headerSvc.setTitle("Loading - Content Invalidation Jobs");
		this.now = new Date();
		const idParam = this.route.snapshot.paramMap.get("id");
		if (!idParam) {
			console.error("Missing route 'id' parameter");
			return;
		}
		this.dsID = parseInt(idParam, 10);
		this.jobAPI.getInvalidationJobs({dsID: this.dsID}).then(
			r => {
				this.jobs = r;
			}
		);
		this.dsAPI.getDeliveryServices(this.dsID).then(
			r => {
				this.deliveryservice = r;
				this.headerSvc.setTitle(`${this.deliveryservice.displayName} - Content Invalidation Jobs`);
			}
		);
	}

	/**
	 * Gets whether or not a Job is in-progress.
	 *
	 * @param j The Job to check.
	 * @returns Whether or not `j` is currently in-progress.
	 */
	public isInProgress(j: InvalidationJob): boolean {
		return j.startTime <= this.now && this.endDate(j) >= this.now;
	}

	/**
	 * Handles a click on a
	 *
	 * @param j The ID of the Job to delete.
	 */
	public async deleteJob(j: number): Promise<void> {
		await this.jobAPI.deleteInvalidationJob(j);
		this.jobs = await this.jobAPI.getInvalidationJobs();
	}

	/**
	 * Gets the ending date and time for a content invalidation job.
	 *
	 * @param j The job from which to extract an end date.
	 * @returns The date at which the Job will stop being in effect.
	 */
	public endDate(j: InvalidationJob): Date {
		if (!j.parameters) {
			throw new Error("cannot get end date for job with no parameters");
		}
		const tmp = j.parameters.replace(/h$/, "").split(":");
		if (tmp.length !== 2) {
			throw new Error(`Malformed job parameters: "${j.parameters}" (id: ${j.id})`);
		}
		const ttl = parseInt(tmp[1], 10);
		if (isNaN(ttl)) {
			throw new Error(`Invalid TTL: "${tmp[1]}" (job id: ${j.id})`);
		}
		// I don't know why this is necessary, because Date.getTime *says* it
		// returns a number, but if you take away the type of `start` here, it
		// fails to compile.
		const start: number = j.startTime.getTime();
		return new Date(start + ttl * 60 * 60 * 1000);
	}

	/**
	 * Creates a new job.
	 *
	 * @param e The DOM event that triggered the creation.
	 */
	public newJob(): void {
		const dialogRef = this.dialog.open(NewInvalidationJobDialogComponent, {data: {dsID: this.dsID}});
		dialogRef.afterClosed().subscribe(
			(created) => {
				if (created) {
					this.jobAPI.getInvalidationJobs({dsID: this.dsID}).then(
						resp => {
							this.jobs = resp;
						}
					);
				}
			}
		);
	}

	/**
	 * Handles a user clicking on a Job's "edit" button by opening the edit
	 * dialog.
	 *
	 * @param job The Job to be edited.
	 */
	public editJob(job: InvalidationJob): void {
		const dialogRef = this.dialog.open(NewInvalidationJobDialogComponent, {data: {dsID: this.dsID, job}});
		dialogRef.afterClosed().subscribe(
			created => {
				if (created) {
					this.jobAPI.getInvalidationJobs({dsID: this.dsID}).then(
						resp => this.jobs = resp
					);
				}
			}
		);
	}
}
