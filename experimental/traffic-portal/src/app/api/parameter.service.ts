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
import {RequestParameter, ResponseParameter} from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * ParameterService exposes API functionality related to Parameters.
 */
@Injectable()
export class ParameterService extends APIService {

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}

	public async getParameters(idOrName: number | string): Promise<ResponseParameter>;
	public async getParameters(): Promise<Array<ResponseParameter>>;
	/**
	 * Retrieves Parameters from the API.
	 *
	 * @param idOrName Specify either the integral, unique identifier (number) of a specific Parameter to retrieve, or its name (string).
	 * @returns The requested Parameter(s).
	 */
	public async getParameters(idOrName?: number | string): Promise<Array<ResponseParameter> | ResponseParameter> {
		const path = "parameters";
		if (idOrName !== undefined) {
			let params;
			switch (typeof idOrName) {
				case "number":
					params = {id: idOrName};
					break;
				case "string":
					params = {name: idOrName};
			}
			const r = await this.get<[ResponseParameter]>(path, undefined, params).toPromise();
			if (r.length !== 1) {
				throw new Error(`Traffic Ops responded with ${r.length} Types by identifier ${idOrName}`);
			}
			return r[0];
		}
		return this.get<Array<ResponseParameter>>(path).toPromise();
	}

	/**
	 * Deletes an existing parameter.
	 *
	 * @param typeOrId Id of the parameter to delete.
	 * @returns The deleted parameter.
	 */
	public async deleteParameter(typeOrId: number | ResponseParameter) {
		const id = typeof(typeOrId) === "number" ? typeOrId : typeOrId.id;
		return this.delete<ResponseParameter>(`parameters/${id}`);
	}

	/**
	 * Creates a new parameter.
	 *
	 * @param parameter The parameter to create.
	 * @returns The created parameter.
	 */
	public async createParameter(parameter: RequestParameter): Promise<ResponseParameter> {
		return this.post<ResponseParameter>("parameters", parameter).toPromise();
	}

	/**
	 * Replaces the current definition of a parameter with the one given.
	 *
	 * @param parameter The new parameter.
	 * @returns The updated parameter.
	 */
	public async updateParameter(parameter: ResponseParameter): Promise<ResponseParameter> {
		const path = `parameters/${parameter.id}`;
		return this.put<ResponseParameter>(path, parameter).toPromise();
	}
}
