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

import { Injectable } from "@angular/core";
import { RequestInvalidationJob, ResponseDeliveryService, ResponseInvalidationJob, ResponseUser } from "trafficops-types";

// This needs to be imported from above, because that's how the services are
// specified in `providers`.
import { DeliveryServiceService } from "..";

/**
 * JobOpts defines the options that can be passed to getInvalidationJobs.
 */
interface JobOpts {
	/** return only the Jobs that operate on this Delivery Service */
	deliveryService?: ResponseDeliveryService;
	/** return only the Jobs that operate on the Delivery Service with this ID */
	dsID?: number;
	/** return only the Job that has this ID */
	id?: number;
	/** return only the Jobs that were created by this user */
	user?: ResponseUser;
	/** return only the Jobs that were created by the user that has this ID */
	userId?: number;
}

/**
 * InvalidationJobService exposes API functionality related to Content Invalidation Jobs.
 */
@Injectable()
export class InvalidationJobService {

	private readonly jobs = new Array<ResponseInvalidationJob>();
	private idCounter = 0;

	constructor(private readonly dsService: DeliveryServiceService) {}

	/**
	 * Fetches all invalidation jobs that match the passed criteria.
	 *
	 * @param opts Optional identifiers for the requested Jobs.
	 * @returns The request Jobs.
	 */
	public async getInvalidationJobs(opts?: JobOpts): Promise<Array<ResponseInvalidationJob>> {
		let ret = this.jobs;
		if (opts) {
			if (opts.id) {
				ret = ret.filter(j=>j.id===opts.id);
			}
			if (opts.dsID) {
				const ds = await this.dsService.getDeliveryServices(opts.dsID);
				ret = ret.filter(j=>j.deliveryService === ds.xmlId);
			}
			if (opts.userId) {
				// TODO: implement this
				throw new Error("filtering by userId not implemented in testing services");
			}
			if (opts.deliveryService && opts.deliveryService.xmlId) {
				ret = ret.filter(j=>j.deliveryService === opts.deliveryService?.xmlId);
			}
			if (opts.user) {
				ret = ret.filter(j=>j.createdBy === opts.user?.username);
			}
		}
		return ret;
	}

	/**
	 * Creates an Invalidation Job.
	 *
	 * @param job The Job to create.
	 * @returns whether or not creation succeeded.
	 */
	public async createInvalidationJob(job: RequestInvalidationJob): Promise<ResponseInvalidationJob> {
		const ret = {
			// Yes, this is ill-formed.
			assetUrl: job.regex,
			createdBy: "test-admin",
			deliveryService: job.deliveryService,
			id: ++this.idCounter,
			invalidationType: job.invalidationType,
			startTime: job.startTime instanceof Date ? job.startTime : new Date(job.startTime),
			ttlHours: job.ttlHours
		};
		this.jobs.push(ret);
		return ret;
	}

	/**
	 * Updates a Job by replacing it with a new definition.
	 *
	 * @param job The new definition of the Job.
	 * @returns The edited Job as returned by the server.
	 */
	public async updateInvalidationJob(job: ResponseInvalidationJob): Promise<ResponseInvalidationJob> {
		const idx = this.jobs.findIndex(j=>j.id===job.id);
		if (idx < 0) {
			throw new Error(`no such Job: #${job.id}`);
		}
		this.jobs[idx] = job;
		return job;
	}

	/**
	 * Deletes a Job.
	 *
	 * @param id The ID of the Job to delete.
	 * @returns The deleted Job.
	 */
	public async deleteInvalidationJob(id: number): Promise<ResponseInvalidationJob> {
		const idx = this.jobs.findIndex(j=>j.id===id);
		if (idx < 0) {
			throw new Error(`no such Job: #${id}`);
		}
		return this.jobs.splice(idx, 1)[0];
	}
}
