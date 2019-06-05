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

export enum JobType {
	PURGE
}

export interface InvalidationJob {
	assetURL?: RegExp;
	createdBy?: string;
	deliveryService?: string;
	dsId?: number;
	id?: number;
	keyword?: string;

	/**
	 * though not enforced by the API (or database), this should ALWAYS have the format 'TTL:nh',
	 * describing the job's TTL in hours (`n` can be any integer value > 0).
	**/
	parameters?: string;
	regex?: RegExp | string;
	startTime: Date;
	ttl?: number;
}
