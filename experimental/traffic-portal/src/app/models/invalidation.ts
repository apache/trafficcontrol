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
/** JobType enumerates the valid types of Job. */
export const enum JobType {
	/** Content Invalidation Request. */
	PURGE = "PURGE"
}

/**
 * InvalidationJob objects represent periods of time over which specific objects
 * may not be cached.
 */
export interface InvalidationJob {
	/**
	 * A regular expression that matches content to be "invalidated" or
	 * "revalidated".
	 */
	assetUrl: string;
	/**
	 * The name of the user that created the Job.
	 */
	createdBy: string;
	/** The XMLID of the Delivery Service within which the Job will operate. */
	deliveryService: string;
	/** An integral, unique identifier for this Job. */
	readonly id: number;
	/** The type of Job. */
	keyword: JobType;

	/**
	 * though not enforced by the API (or database), this should ALWAYS have the
	 * format 'TTL:nh', describing the job's TTL in hours (`n` can be any
	 * integer value > 0).
	 */
	parameters: string;
	/**
	 * The time at which the Job is scheduled to start.
	 */
	startTime: Date;
}
