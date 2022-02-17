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

import { HttpClient, HttpHeaders } from "@angular/common/http";
import type { Observable } from "rxjs";
import { map } from "rxjs/operators";

import { environment } from "src/environments/environment";

/**
 * This is the base class from which all other API classes inherit.
 */
export abstract class APIService {
	/**
	 * The API version used by the service(s) - this will be overridden by the
	 * environment if a different API version is therein found.
	 */
	public apiVersion = "3.0";

	/**
	 * Sends an HTTP DELETE request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	protected delete<T>(path: string, data?: object, params?: Record<string, string>): Observable<T> {
		return this.do<T>("delete", path, data, params);
	}

	/**
	 * Sends an HTTP GET request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	protected get<T>(path: string, data?: object, params?: Record<string, string>): Observable<T> {
		return this.do<T>("get", path, data, params);
	}

	/**
	 * Sends an HTTP HEAD request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	protected head<T>(path: string, data?: object, params?: Record<string, string>): Observable<T> {
		return this.do<T>("head", path, data, params);
	}

	/**
	 * Sends an HTTP OPTIONS request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	protected options<T>(path: string, data?: object, params?: Record<string, string>): Observable<T> {
		return this.do<T>("options", path, data, params);
	}

	/**
	 * Sends an HTTP PATCH request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	protected patch<T>(path: string, data?: object, params?: Record<string, string>): Observable<T> {
		return this.do<T>("patch", path, data, params);
	}

	/**
	 * Sends an HTTP POST request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	protected post<T>(path: string, data?: object, params?: Record<string, string>): Observable<T> {
		return this.do<T>("post", path, data, params);
	}

	/**
	 * Sends an HTTP PUT request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	protected put<T>(path: string, data?: object, params?: Record<string, string>): Observable<T> {
		return this.do<T>("put", path, data, params);
	}

	/**
	 * Sends an HTTP request to the API.
	 *
	 * @param method The HTTP request method to use, e.g. "GET".
	 * @param path The request path.
	 * @param body Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	protected do<T>(method: string, path: string, body?: object, params?: Record<string, string>): Observable<T> {

		const options = {
			body,
			params,
			...this.defaultOptions
		};
		// TODO pass alerts to the alert service
		// (TODO create the alert service)
		return this.http.request<{response: T}>(method, `/api/${this.apiVersion}/${path.replace(/^\/+/, "")}`, options).pipe(map(
			r => {
				if (!r.body) {
					throw new Error(`${method} ${path} returned no response body - ${r.status} ${r.statusText}`);
				}
				return r.body.response;
			}
		));
	}

	/**
	 * These are the default options sent in HttpClient methods - if subclasses
	 * make raw requests, they are encouraged to extend this rather than duplicate
	 * it.
	 */
	protected readonly defaultOptions = {
		// This is part of the HTTP spec. I can't - and shouldn't - change it.
		// eslint-disable-next-line @typescript-eslint/naming-convention
		headers: new HttpHeaders({"Content-Type": "application/json"}),
		observe: "response" as "response",
		responseType: "json" as "json",
	};

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
