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
import { ResponseDeliveryService, ResponseInvalidationJob } from "trafficops-types";

import { DeliveryServiceService, InvalidationJobService } from "src/app/api";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import {
	NewInvalidationJobDialogComponent,
	type NewInvalidationJobDialogData
} from "./new-invalidation-job-dialog/new-invalidation-job-dialog.component";

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
	public deliveryservice!: ResponseDeliveryService;

	/** All of the jobs for the described Delivery Service. */
	public jobs: Array<ResponseInvalidationJob>;

	/** The current date/time when the page loads */
	public now: Date = new Date();

	/** The ID of the Delivery Service to which the displayed Jobs belong. */
	private dsID = -1;
	constructor(
		private readonly route: ActivatedRoute,
		private readonly jobAPI: InvalidationJobService,
		private readonly dsAPI: DeliveryServiceService,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) {
		this.jobs = new Array<ResponseInvalidationJob>();
	}

	/**
	 * Runs initialization, fetching the jobs and Delivery Service data from
	 * Traffic Ops and setting the date/time on page load.
	 */
	public async ngOnInit(): Promise<void> {
		this.navSvc.headerTitle.next("Loading - Content Invalidation Jobs");
		this.now = new Date();
		const idParam = this.route.snapshot.paramMap.get("id");
		if (!idParam) {
			this.log.error("Missing route 'id' parameter");
			return;
		}
		this.dsID = parseInt(idParam, 10);
		this.jobs = await this.jobAPI.getInvalidationJobs({dsID: this.dsID});
		this.deliveryservice = await this.dsAPI.getDeliveryServices(this.dsID);
		this.navSvc.headerTitle.next(`${this.deliveryservice.displayName} - Content Invalidation Jobs`);
	}

	/**
	 * Gets whether or not a Job is in-progress.
	 *
	 * @param j The Job to check.
	 * @returns Whether or not `j` is currently in-progress.
	 */
	public isInProgress(j: ResponseInvalidationJob): boolean {
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
	public endDate(j: ResponseInvalidationJob): Date {
		const start = j.startTime.getTime();
		return new Date(start + j.ttlHours * 60 * 60 * 1000);
	}

	/**
	 * Creates a new job.
	 *
	 * @param e The DOM event that triggered the creation.
	 */
	public newJob(): void {
		const data: NewInvalidationJobDialogData = {
			dsID: this.deliveryservice.xmlId
		};
		const dialogRef = this.dialog.open(NewInvalidationJobDialogComponent, {data});
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
	public editJob(job: ResponseInvalidationJob): void {
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
