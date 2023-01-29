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

import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import type { RequestInvalidationJob, ResponseDeliveryService, ResponseInvalidationJob, ResponseUser } from "trafficops-types";

import { APIService } from "./base-api.service";

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
export class InvalidationJobService extends APIService {

	constructor(http: HttpClient) {
		super(http);
	}

	/**
	 * Fetches all invalidation jobs that match the passed criteria.
	 *
	 * @param opts Optional identifiers for the requested Jobs.
	 * @returns The request Jobs.
	 */
	public async getInvalidationJobs(opts?: JobOpts): Promise<Array<ResponseInvalidationJob>> {
		const path = "jobs";
		const params: Record<string, string> = {};
		if (opts) {
			if (opts.id) {
				params.id = String(opts.id);
			}
			if (opts.dsID) {
				params.dsId = String(opts.dsID);
			}
			if (opts.userId) {
				params.userId = String(opts.userId);
			}
			if (opts.deliveryService && opts.deliveryService.id) {
				params.dsId = String(opts.deliveryService.id);
			}
			if (opts.user && opts.user.id) {
				params.userId = String(opts.user.id);
			}
		}
		const js = await this.get<Array<ResponseInvalidationJob>>(path, undefined, params).toPromise();
		const jobs = new Array<ResponseInvalidationJob>();
		for (const j of js) {
			const tmp = String(j.startTime).replace(" ", "T").replace("+00", "Z");
			j.startTime = new Date(tmp);
			jobs.push(j);
		}
		return jobs;
	}

	/**
	 * Creates an Invalidation Job.
	 *
	 * @param job The Job to create.
	 * @returns whether or not creation succeeded.
	 */
	public async createInvalidationJob(job: RequestInvalidationJob): Promise<ResponseInvalidationJob> {
		const path = "jobs";
		return this.post<ResponseInvalidationJob>(path, job).toPromise();
	}

	/**
	 * Updates a Job by replacing it with a new definition.
	 *
	 * @param job The new definition of the Job.
	 * @returns The edited Job as returned by the server.
	 */
	public async updateInvalidationJob(job: ResponseInvalidationJob): Promise<ResponseInvalidationJob> {
		return this.put<ResponseInvalidationJob>("jobs", job, {id: String(job.id)}).toPromise();
	}

	/**
	 * Deletes a Job.
	 *
	 * @param id The ID of the Job to delete.
	 * @returns The deleted Job.
	 */
	public async deleteInvalidationJob(id: number): Promise<ResponseInvalidationJob> {
		return this.delete<ResponseInvalidationJob>("jobs", undefined, {id: String(id)}).toPromise();
	}
}
