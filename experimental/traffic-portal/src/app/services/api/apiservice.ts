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

	/**
	 * Sends an HTTP DELETE request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @returns An Observable that emits the server response.
	 */
	protected delete(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("delete", path, data);
	}

	/**
	 * Sends an HTTP GET request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @returns An Observable that emits the server response.
	 */
	protected get(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("get", path, data);
	}

	/**
	 * Sends an HTTP HEAD request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @returns An Observable that emits the server response.
	 */
	protected head(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("head", path, data);
	}

	/**
	 * Sends an HTTP OPTIONS request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @returns An Observable that emits the server response.
	 */
	protected options(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("options", path, data);
	}

	/**
	 * Sends an HTTP PATCH request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @returns An Observable that emits the server response.
	 */
	protected patch(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("patch", path, data);
	}

	/**
	 * Sends an HTTP POST request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @returns An Observable that emits the server response.
	 */
	protected post(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("post", path, data);
	}

	/**
	 * Sends an HTTP PUSH request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @returns An Observable that emits the server response.
	 */
	protected push(path: string, data?: object): Observable<HttpResponse<object>> {
		return this.do("push", path, data);
	}

	/**
	 * Sends an HTTP request to the API.
	 *
	 * @param method The HTTP request method to use, e.g. "GET".
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @returns An Observable that emits the server response.
	 */
	protected do(method: string, path: string, data?: object): Observable<HttpResponse<object>> {

		const options = {
			body: data,
			// This is part of the HTTP spec. I can't - and shouldn't - change it.
			// eslint-disable-next-line @typescript-eslint/naming-convention
			headers: new HttpHeaders({"Content-Type": "application/json"}),
			observe: "response" as "response",
			responseType: "json" as "json",
		};
		// TODO pass alerts to the alert service
		// (TODO create the alert service)
		return this.http.request(method, path, options);
	}

	/**
	 * Constructs the service and sets the API version based on the execution environment.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(protected readonly http: HttpClient) {
		if (environment.apiVersion) {
			this.apiVersion = environment.apiVersion;
		}
	}
}
