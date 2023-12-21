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

import { type HttpClient, HttpHeaders, HttpRequest, HttpParams } from "@angular/common/http";
import type { Observable } from "rxjs";

import { environment } from "src/environments/environment";

import type { APIService, QueryParams } from "../base-api.service";

/**
 * This is the base class from which all other API classes inherit.
 */
export abstract class APITestingService implements APIService {
	/**
	 * The API version used by the service(s) - this will be overridden by the
	 * environment if a different API version is therein found.
	 */
	public readonly apiVersion = environment.apiVersion;

	/**
	 * This exists to satisfy typing requirements, but always has the actual,
	 * underlying value of `undefined`, so **do not use it**.
	 */
	public http!: HttpClient;

	/**
	 * Holds a stack of the requests the API service has "sent". Most tests
	 * probably won't need this, but since we need to implement `do` anyway I
	 * figured this was an easy, useful way to satisfy typings.
	 */
	public requestStack = new Array<HttpRequest<unknown>>();

	/**
	 * These are the default options sent in HttpClient methods - if subclasses
	 * make raw requests, they are encouraged to extend this rather than duplicate
	 * it.
	 */
	public readonly defaultOptions = {
		// This is part of the HTTP spec. I can't - and shouldn't - change it.
		// eslint-disable-next-line @typescript-eslint/naming-convention
		headers: new HttpHeaders({"Content-Type": "application/json"}),
		observe: "response" as "response",
		responseType: "json" as "json",
	};

	/**
	 * Sends an HTTP DELETE request to the API.
	 *
	 * @param path The request path.
	 * @param data Optional request body (will be JSON.stringify'd).
	 * @param params Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	public delete<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	public get<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	public head<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	public options<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	public patch<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	public post<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
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
	public put<T = undefined>(path: string, data?: object, params?: QueryParams): Observable<T> {
		return this.do<T>("put", path, data, params);
	}

	/**
	 * Sends an HTTP request to the API. The testing implementation should be
	 * treated as having the concrete return type `void`, because **it does not
	 * return any value**. Instead, it pushes a request to the
	 * {@link APITestingService.requestStack}.
	 *
	 * @param method The HTTP request method to use, e.g. "GET".
	 * @param path The request path.
	 * @param body Optional request body (will be JSON.stringify'd).
	 * @param qParams Option query parameters to send in the request.
	 * @returns An Observable that emits the server response.
	 */
	public do<T>(method: string, path: string, body?: object, qParams?: QueryParams): Observable<T> {

		const params = qParams instanceof HttpParams ? qParams : new HttpParams({fromObject: qParams});

		const options = {
			body,
			params,
			...this.defaultOptions
		};

		switch(method.toUpperCase()) {
			case "GET":
			case "HEAD":
			case "OPTIONS":
			case "DELETE":
				this.requestStack.push(new HttpRequest(method, `/api/${this.apiVersion}/${path.replace(/^\/+/, "")}`, options));
				break;
			case "PUT":
			case "POST":
			case "PATCH":
				this.requestStack.push(new HttpRequest(method, `/api/${this.apiVersion}/${path.replace(/^\/+/, "")}`, body, options));
				break;
		}

		return undefined as unknown as Observable<T>;
	}
}
