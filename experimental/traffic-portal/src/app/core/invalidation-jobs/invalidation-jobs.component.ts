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
import { ActivatedRoute } from "@angular/router";
import { faPlus } from "@fortawesome/free-solid-svg-icons";

import { defaultDeliveryService, DeliveryService, InvalidationJob } from "../../models";
import { DeliveryServiceService, InvalidationJobService } from "../../shared/api";
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
	private dsId = -1;

	/** The icon for the "Create a new Job" FAB. */
	public readonly addIcon = faPlus;

	/**
	 * Constructor.
	 */
	constructor(
		private readonly route: ActivatedRoute,
		private readonly jobAPI: InvalidationJobService,
		private readonly dsAPI: DeliveryServiceService,
		private readonly dialog: MatDialog
	) {
		this.deliveryservice = {...defaultDeliveryService};
		this.jobs = new Array<InvalidationJob>();
	}

	/**
	 * Runs initialization, fetching the jobs and Delivery Service data from
	 * Traffic Ops and setting the pageload date/time.
	 */
	public ngOnInit(): void {
		this.now = new Date();
		const idParam = this.route.snapshot.paramMap.get("id");
		if (!idParam) {
			console.error("Missing route 'id' parameter");
			return;
		}
		this.dsId = parseInt(idParam, 10);
		this.jobAPI.getInvalidationJobs({dsID: this.dsId}).then(
			r => {
				// The values returned by the API are not RFC-compliant at the time of this writing,
				// so we need to do some pre-processing on them.
				for (const j of r) {
					const tmp = Array.from(String(j.startTime).split(" ").join("T"));
					tmp.splice(-3, 3);
					j.startTime = new Date(tmp.join(""));
					this.jobs.push(j);
				}
			}
		);
		this.dsAPI.getDeliveryServices(this.dsId).then(
			r => {
				this.deliveryservice = r;
			}
		);
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
		const tmp = j.parameters.split(":");
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
		const dialogRef = this.dialog.open(NewInvalidationJobDialogComponent, {data: this.dsId});
		dialogRef.afterClosed().subscribe(
			(created) => {
				if (created) {
					this.jobAPI.getInvalidationJobs({dsID: this.dsId}).then(
						resp => {
							this.jobs = new Array<InvalidationJob>();
							for (const j of resp) {
								const tmp = Array.from(String(j.startTime).replace(" ", "T"));
								tmp.splice(-3, 3);
								j.startTime = new Date(tmp.join(""));
								this.jobs.push(j);
							}
						}
					);
				}
			}
		);
	}
}
