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
import { map } from "rxjs/operators";

import { environment } from "../../../environments/environment";

/**
 * This is the base class from which all other API classes inherit.
 */
export class APIService {
	/**
	 * The API version used by the service(s) - this will be overridden by the
	 * environment if a different API version is therein found.
	 */
	public API_VERSION = "2.0";

	protected delete (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do("delete", path, data);
	}
	protected get (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do("get", path, data);
	}
	protected head (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do("head", path, data);
	}
	protected options (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do("options", path, data);
	}
	protected patch (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do("patch", path, data);
	}
	protected post (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do("post", path, data);
	}
	protected push (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do("push", path, data);
	}

	protected do (method: string, path: string, data?: Object): Observable<HttpResponse<any>> {

		/* tslint:disable */
		const options = {headers: new HttpHeaders({'Content-Type': 'application/json'}),
		                 observe: 'response' as 'response',
		                 responseType: 'json' as 'json',
		                 body: data};
		/* tslint:enable */
		return this.http.request(method, path, options).pipe(map((response) => {
			// TODO pass alerts to the alert service
			// (TODO create the alert service)
			return response as HttpResponse<any>;
		}));
	}

	constructor(private readonly http: HttpClient) {
		if (environment.APIVersion) {
			this.API_VERSION = environment.APIVersion;
		}
	}
}
