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
import { HttpRequest, HttpHandler, HttpEvent, HttpInterceptor } from "@angular/common/http";
import { Injectable } from "@angular/core";

import { Observable, throwError } from "rxjs";
import { catchError } from "rxjs/operators";

import { Alert } from "../models/alert";
import { AlertService, CurrentUserService } from "../services";

/**
 * This class intercepts any and all HTTP error responses and checks for
 * authorization problems. It then redirects the user back to login.
 */
@Injectable()
export class ErrorInterceptor implements HttpInterceptor {

	constructor(
		private readonly currentUserService: CurrentUserService,
		private readonly alerts: AlertService
	) {}

	/**
	 * Intercepts HTTP responses and checks for erroneous responses, displaying
	 * appropriate error Alerts and redirecting unauthenticated users to the
	 * login form.
	 *
	 * @param request The client request.
	 * @param next The next handler for HTTP requests in the pipeline.
	 * @returns An Observable that will emit an event if the request fails.
	 */
	public intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
		return next.handle(request).pipe(catchError((err) => {
			console.error("HTTP Error: ", err);

			if (err.hasOwnProperty("error") && (err as {error: object}).error.hasOwnProperty("alerts")) {
				for (const a of (err as {error: {alerts: Alert[]}}).error.alerts) {
					this.alerts.alertsSubject.next(a);
				}
			}
			if (err.status === 401 || err.status === 403) {
				this.currentUserService.logout(true);
			}

			return throwError(err);
		}));
	}
}
