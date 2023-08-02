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
import { type HttpRequest, type HttpHandler, type HttpEvent, type HttpInterceptor, HttpErrorResponse } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Router } from "@angular/router";
import { type Observable, throwError } from "rxjs";
import { catchError } from "rxjs/operators";
import type { Alert } from "trafficops-types";

import { AlertService } from "../alert/alert.service";
import { LoggingService } from "../logging.service";

/**
 * This class intercepts any and all HTTP error responses and checks for
 * authorization problems. It then redirects the user back to login.
 */
@Injectable()
export class ErrorInterceptor implements HttpInterceptor {

	constructor(
		private readonly alerts: AlertService,
		private readonly router: Router,
		private readonly log: LoggingService,
	) {}

	/**
	 * Raises all passed alerts to the AlertService.
	 *
	 * @param alerts The alerts to be raised.
	 */
	private raiseAlerts(alerts: Alert[]): void {
		for (const alert of alerts) {
			this.alerts.newAlert(alert);
		}
	}

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
		return next.handle(request).pipe(catchError((err: HttpErrorResponse) => {
			// I don't know why, but sometimes these errors have just no content
			// and stringify to simply just the word "Error". So in order to get
			// anything at all useful out of them, I'm adding a stack trace at
			// the debugging level.
			this.log.error(`HTTP error: ${err.message || err.error || err}`);
			this.log.debug(err);

			if (typeof(err.error) === "string") {
				try {
					const body: {alerts: Alert[] | undefined} = JSON.parse(err.error);
					if (Array.isArray(body.alerts)) {
						this.raiseAlerts(body.alerts);
					}
				} catch (e) {
					this.log.error("non-JSON HTTP error response:", e);
				}
			} else if (typeof(err.error) === "object" && Array.isArray(err.error.alerts)) {
				this.raiseAlerts(err.error.alerts);
			}

			if (err instanceof HttpErrorResponse && err.status === 401 && this.router.getCurrentNavigation() === null) {
				const currentURL = this.router.parseUrl(this.router.routerState.snapshot.url);
				const path = Object.entries(currentURL.root.children).map(c=>c[1].segments.map(s=>s.path).join("/")).join("/");
				if (path !== "login") {
					const params = Object.entries(currentURL.queryParams).filter(([param])=>param!=="returnUrl");
					let returnUrl = path;
					if (params.length > 0) {
						returnUrl = `${returnUrl}?${params.map(([k,v])=>`${k}=${v}`).join("&")}`;
					}
					if (currentURL.fragment) {
						returnUrl += `#${currentURL.fragment}`;
					}
					this.router.navigate(["/login"], {queryParams: {returnUrl}});
				}
			}

			return throwError(err);
		}));
	}
}
