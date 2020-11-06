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
import { FormControl } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";

import { Subject } from "rxjs";

import { DeliveryService, InvalidationJob } from "../../models";
import { DeliveryServiceService, InvalidationJobService } from "../../services/api";

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
	public now: Date;

	/** Whether or not to show the dialog for creating a new job. */
	public showDialog: Subject<boolean>;

	private dsId: number;


	/** Control for users to enter new content invalidation jobs. */
	public regexp = new FormControl("/");
	/** Control for users to enter a new job's TTL. */
	public ttl = new FormControl(178);
	/** Control for users to enter the starting date for a new job. */
	public startDate = new FormControl("");
	/** Control for users to enter the starting time for a new job. */
	public startTime = new FormControl("");
	/**
	 * Sets a customvalidity message when the user-entered regular expression is
	 * not valid.
	 */
	public regexpIsValid: Subject<string>;


	constructor (
		private readonly route: ActivatedRoute,
		private readonly jobAPI: InvalidationJobService,
		private readonly dsAPI: DeliveryServiceService
	) {
		this.deliveryservice = {active: true} as DeliveryService;
		this.jobs = new Array<InvalidationJob>();
		this.showDialog = new Subject<boolean>();
		this.regexpIsValid = new Subject<string>();
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
		this.jobAPI.getInvalidationJobs({dsID: this.dsId}).subscribe(
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
		this.dsAPI.getDeliveryServices(this.dsId).subscribe(
			(r: DeliveryService) => {
				this.deliveryservice = r;
			}
		);
	}

	/** Gets the ending date and time for a content invalidation job. */
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
		return new Date(new Date(j.startTime.getTime() + ttl * 60 * 60 * 1000));
	}

	/**
	 * Creates a new job.
	 *
	 * @param e The DOM event that triggered the creation.
	 */
	public newJob(e?: Event): void {
		if (e) {
			e.preventDefault();
			e.stopPropagation();
		}

		const now = new Date();
		now.setUTCMilliseconds(0);

		const year = String(now.getFullYear()).padStart(4, "0");
		const month = String(now.getMonth()+1).padStart(2, "0");
		const day = String(now.getDate()).padStart(2, "0");
		this.startDate.setValue(`${year}-${month}-${day}`);

		const hours = String(now.getHours()).padStart(2, "0");
		const minutes = String(now.getMinutes()).padStart(2, "0");
		this.startTime.setValue(`${hours}:${minutes}`);

		this.showDialog.next(true);
	}

	/**
	 * Closes the "new content invalidation job" dialog, canceling the creation.
	 */
	public closeDialog(e: Event): void {
		e.preventDefault();
		e.stopPropagation();
		this.showDialog.next(false);
	}

	/**
	 * Closes the "new content invalidation job" dialog, creating the new job.
	 */
	public submitDialog(e: Event): void {
		e.preventDefault();
		e.stopPropagation();

		let re: RegExp;
		try {
			re = new RegExp(this.regexp.value);
		} catch (err) {
			this.regexpIsValid.next(`Must be a valid regular expression! (${err})`);
			return;
		}

		const job = {
			dsId: this.deliveryservice.id,
			parameters: `TTL:${this.ttl.value}`,
			regex: re.toString().replace(/^\/|\/$/g, "").replace("\\/", "/"),
			startTime: this.startDate.value.concat(" ", `${this.startTime.value}:00`),
			ttl: this.ttl.value
		};

		this.jobAPI.createInvalidationJob(job).subscribe(
			r => {
				if (r) {
					this.jobAPI.getInvalidationJobs({dsID: this.dsId}).subscribe(
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
					this.showDialog.next(false);
				} else {
					console.warn("failure");
				}
			},
			err => {
				console.error("error: ", err);
			}
		);
	}

}
