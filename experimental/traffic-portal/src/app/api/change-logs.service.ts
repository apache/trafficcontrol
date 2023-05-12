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
import { Log } from "trafficops-types";

import { APIService, type QueryParams } from "src/app/api/base-api.service";

/**
 * Defines & handles api endpoints related to change logs
 */
@Injectable()
export class ChangeLogsService extends APIService {

	constructor(http: HttpClient) {
		super(http);
	}

	/**
	 * Calls api logs endpoint
	 *
	 * @param params Request parameters to add
	 * @returns Change logs
	 */
	public async getChangeLogs(params?: QueryParams): Promise<Array<Log>> {
		return this.get<Array<Log>>("logs", undefined, params).toPromise();
	}
}
