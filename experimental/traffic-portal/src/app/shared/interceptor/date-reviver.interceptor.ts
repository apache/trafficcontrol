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

import {
	type HttpRequest,
	type HttpHandler,
	type HttpEvent,
	type HttpInterceptor,
	HttpResponse
} from "@angular/common/http";
import { Injectable } from "@angular/core";
import type { Observable } from "rxjs";
import { map } from "rxjs/operators";

import { dateReviver } from "src/app/utils/date";

/**
 * The DateReviverInterceptor adds custom JSON parsing to all HTTP requests that
 * converts date strings to Date instances.
 */
@Injectable()
export class DateReviverInterceptor implements HttpInterceptor {

	/**
	 * Parses the eventually received response (if indeed one is received) as a
	 * JSON payload, reviving string values that look like dates into Date
	 * instances.
	 *
	 * @param request The outgoing request.
	 * @param next The next step in the request process.
	 * @returns The response events. Non-response events or those without
	 * response bodies are untouched.
	 */
	private parseResponseJSON(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
		request = request.clone({responseType: "text"});
		return next.handle(request).pipe(map(
			event => {
				if (event instanceof HttpResponse && typeof event.body === "string") {
					return event.clone({body: JSON.parse(event.body, dateReviver)});
				}
				return event;
			}));
	}

	/**
	 * Intercepts requests with the "json" response type to add a custom parser
	 * for date strings.
	 *
	 * @param request The outgoing request, before being sent.
	 * @param next The next step in the HTTP handler stack.
	 * @returns The response events. Non-response events or those without
	 * response bodies are untouched.
	 */
	public intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
		if (request.responseType === "json") {
			// If the expected response type is JSON then handle it here.
			return this.parseResponseJSON(request, next);
		}
		return next.handle(request);
	}
}
