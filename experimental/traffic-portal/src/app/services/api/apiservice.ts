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

import { HttpClient, HttpResponse, HttpHeaders } from "@angular/common/http";

import { Observable } from "rxjs";

import { environment } from "../../../environments/environment";

/**
 * This is the base class from which all other API classes inherit.
 */
export class APIService {
	/**
	 * The API version used by the service(s) - this will be overridden by the
	 * environment if a different API version is therein found.
	 */
	public apiVersion = "2.0";

	protected delete(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("delete", path, data);
	}
	protected get(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("get", path, data);
	}
	protected head(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("head", path, data);
	}
	protected options(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("options", path, data);
	}
	protected patch(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("patch", path, data);
	}
	protected post(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("post", path, data);
	}
	protected push(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("push", path, data);
	}

	protected do(method: string, path: string, data?: object): Observable<HttpResponse<object>> {

		/* eslint-disable */
		const options = {headers: new HttpHeaders({'Content-Type': 'application/json'}),
		                 observe: 'response' as 'response',
		                 responseType: 'json' as 'json',
		                 body: data};
		/* eslint-enable */
		// TODO pass alerts to the alert service
		// (TODO create the alert service)
		return this.http.request(method, path, options);
	}

	constructor(private readonly http: HttpClient) {
		if (environment.apiVersion) {
			this.apiVersion = environment.apiVersion;
		}
	}
}
