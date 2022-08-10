/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { Injectable } from "@angular/core";
import { Log } from "trafficops-types";

/**
 * Defines & handles api endpoints related to change logs
 */
@Injectable()
export class ChangeLogsService {

	private readonly changeLogs: Array<Log> = [{
		id: 0,
		lastUpdated: new Date(),
		level: "APICHANGE",
		message: "msg 1",
		ticketNum: null,
		user: "user 1"
	},
	{
		id: 1,
		lastUpdated: new Date(),
		level: "APICHANGE",
		message: "msg 2",
		ticketNum: null,
		user: "user 2"
	},
	{
		id: 2,
		lastUpdated: new Date(),
		level: "APICHANGE",
		message: "msg 3",
		ticketNum: null,
		user: "user 3"
	}
	];

	/**
	 * Calls api logs endpoint
	 *
	 * @param params Request parameters to add
	 * @returns Change logs
	 */
	public async getChangeLogs(params?: Record<string, string>): Promise<Array<Log>> {
		if (params === undefined) {
			return this.changeLogs;
		}
		if("user" in params) {
			return this.changeLogs.filter(cl => cl.user === params.user);
		}
		if("days" in params) {
			return this.changeLogs;
		}
		throw new Error(`unknown params ${params}`);
	}
}
