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

import { DeliveryService, InvalidationJob, NewInvalidationJob, User } from "../../models";

import { APIService } from "./APIService";

/**
 * JobOpts defines the options that can be passed to getInvalidationJobs.
 */
interface JobOpts {
	/** return only the Jobs that operate on this Delivery Service */
	deliveryService?: DeliveryService;
	/** return only the Jobs that operate on the Delivery Service with this ID */
	dsID?: number;
	/** return only the Job that has this ID */
	id?: number;
	/** return only the Jobs that were created by this user */
	user?: User;
	/** return only the Jobs that were created by the user that has this ID */
	userId?: number;
}

/**
 * InvalidationJobService exposes API functionality related to Content Invalidation Jobs.
 */
@Injectable()
export class InvalidationJobService extends APIService {

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}

	/**
	 * Fetches all invalidation jobs that match the passed criteria.
	 *
	 * @param opts Optional identifiers for the requested Jobs.
	 * @returns The request Jobs.
	 */
	public async getInvalidationJobs(opts?: JobOpts): Promise<Array<InvalidationJob>> {
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
		return this.get<Array<InvalidationJob>>(path, undefined, params).toPromise().catch(
			e => {
				console.error("Failed to get Invalidation Jobs:", e);
				return [];
			}
		);
	}

	/**
	 * Creates an Invalidation Job.
	 *
	 * @param job The Job to create.
	 * @returns whether or not creation succeeded.
	 */
	public async createInvalidationJob(job: NewInvalidationJob): Promise<boolean> {
		const path = "jobs";
		return this.post(path, job).toPromise().then(
			() => true,
			() => false
		);
	}
}
