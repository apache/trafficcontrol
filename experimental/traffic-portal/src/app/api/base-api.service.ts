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

import { HttpClient, HttpHeaders, type HttpParams } from "@angular/common/http";
import type { Observable } from "rxjs";
import { map } from "rxjs/operators";
import type { Alert } from "trafficops-types";

import { environment } from "src/environments/environment";

import { hasProperty, isArray } from "../utils";

/**
 * The type of parameters that can be sent in a query string.
 *
 * We're not using @angular/routing.Params because that's too permissive. It's
 * front-end only, so it allows you to pass arbitrary data, even if it's not
 * really representable in a URL query string.
 */
export type QueryParams = Record<Exclude<PropertyKey, symbol>, string | number | boolean> | HttpParams;

/**
 * Checks if something is an Alert.
 *
 * @param x The thing to check.
 * @returns `true` if `x` is an Alert (or at least close enough), `false`
 * otherwise.
 */
function isAlert(x: unknown): x is Alert {
	if (typeof(x) !== "object" || !x) {
		return false;
	}

	return hasProperty(x, "level", "string") && hasProperty(x, "text", "string");
}

/**
 * Checks if an arbitrary object parsed from a response body is Alerts. This is
 * useful for methods that typically return non-JSON data - except in the event
 * of failures.
 *
 * @param x The object to check.
 * @returns `true` if `x` has an `alerts` array, `false` otherwise.
 */
export function hasAlerts(x: object): x is ({alerts: Alert[]}) {
	if (!hasProperty(x, "alerts")) {
		return false;
	}
	return isArray(x.alerts, isAlert);
}

/**
 * This is the base class from which all other API classes inherit.
 */
export abstract class APIService {
	/**
	 * The API version used by the service(s) - this will be overridden by the
	 * environment if a different API version is therein found.
	 */
	public apiVersion = "4.0";

	/**
	 * Sends an HTTP DELETE request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	protected delete<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	protected get<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	protected head<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	protected options<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	protected patch<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	protected post<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	protected put<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	protected do<T>(method: string, path: string, body?: object, params?: QueryParams): Observable<T> {

		const options = {
			body,
			params,
			...this.defaultOptions
		};
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
