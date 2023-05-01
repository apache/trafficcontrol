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
import {RequestParameter, ResponseParameter} from "trafficops-types";

/**
 * ParameterService exposes API functionality relating to Parameters.
 */
@Injectable()
export class ParameterService {
	private lastID = 20;
	private readonly parameters = [
		{
			configFile: "cfg.txt",
			id: 1,
			lastUpdated: new Date(),
			name: "param1",
			profiles: [],
			secure: false,
			value: "10"
		}
	];

	public async getParameters(idOrName: number | string): Promise<ResponseParameter>;
	public async getParameters(): Promise<Array<ResponseParameter>>;
	/**
	 * Gets one or all Parameters from Traffic Ops
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name (string) of a single parameter to be returned.
	 * @returns The requested parameter(s).
	 */
	public async getParameters(idOrName?: number | string): Promise<ResponseParameter | Array<ResponseParameter>> {
		if (idOrName !== undefined) {
			let parameter;
			switch (typeof idOrName) {
				case "string":
					parameter = this.parameters.filter(t=>t.name === idOrName)[0];
					break;
				case "number":
					parameter = this.parameters.filter(t=>t.id === idOrName)[0];
			}
			if (!parameter) {
				throw new Error(`no such Parameter: ${idOrName}`);
			}
			return parameter;
		}
		return this.parameters;
	}

	/**
	 * Deletes an existing parameter.
	 *
	 * @param id Id of the parameter to delete.
	 * @returns The deleted parameter.
	 */
	public async deleteParameter(id: number): Promise<ResponseParameter> {
		const index = this.parameters.findIndex(t => t.id === id);
		if (index === -1) {
			throw new Error(`no such Parameter: ${id}`);
		}
		return this.parameters.splice(index, 1)[0];
	}

	/**
	 * Creates a new parameter.
	 *
	 * @param parameter The parameter to create.
	 * @returns The created parameter.
	 */
	public async createParameter(parameter: RequestParameter): Promise<ResponseParameter> {
		const t = {
			...parameter,
			configFile: "cfg.txt",
			id: ++this.lastID,
			lastUpdated: new Date(),
			profiles: [],
			secure: false,
			value: "100"
		};
		this.parameters.push(t);
		return t;
	}
}
