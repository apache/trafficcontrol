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

import { Observable } from "rxjs";
import { map } from "rxjs/operators";

import { APIService } from "./apiservice";

import { DeliveryService, InvalidationJob, User } from "../../models";

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
@Injectable({providedIn: "root"})
export class InvalidationJobService extends APIService {

	/**
	 * Fetches all invalidation jobs that match the passed criteria.
	 */
	public getInvalidationJobs (opts?: JobOpts): Observable<Array<InvalidationJob>> {
		let path = `/api/${this.API_VERSION}/jobs`;
		if (opts) {
			const args = new Array<string>();
			if (opts.id) {
				args.push(`id=${opts.id}`);
			}
			if (opts.dsID) {
				args.push(`dsId=${opts.dsID}`);
			}
			if (opts.userId) {
				args.push(`userId=${opts.userId}`);
			}
			if (opts.deliveryService && opts.deliveryService.id) {
				args.push(`dsId=${opts.deliveryService.id}`);
			}
			if (opts.user && opts.user.id) {
				args.push(`userId=${opts.user.id}`);
			}

			if (args.length > 0) {
				path += `?${args.join("&")}`;
			}
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<InvalidationJob>;
			}
		));
	}

	/**
	 * Creates the passed invalidation job, returning whether or not the creation was successful.
	 */
	public createInvalidationJob(job: InvalidationJob): Observable<boolean> {
		const path = `/api/${this.API_VERSION}/user/current/jobs`;
		return this.post(path, job).pipe(map(
			_ => true,
			_ => false
		));
	}

	constructor(http: HttpClient) {
		super(http);
	}
}
