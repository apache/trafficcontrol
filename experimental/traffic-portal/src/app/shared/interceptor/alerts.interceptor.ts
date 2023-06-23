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
import type { HttpRequest, HttpHandler, HttpEvent, HttpInterceptor } from "@angular/common/http";
import { Injectable } from "@angular/core";
import type { Observable } from "rxjs";
import { tap } from "rxjs/operators";
import type { Alert } from "trafficops-types";

import { AlertService } from "../alert/alert.service";

/**
 * This class intercepts any and all alerts contained in API responses and
 * passes them to the `AlertService` for display to the user.
 */
@Injectable()
export class AlertInterceptor implements HttpInterceptor {
	/**
	 * Constructor.
	 */
	constructor(private readonly alertService: AlertService) {
	}

	/**
	 * Intercepts HTTP responses and checks for any alerts.
	 *
	 * @param request The client request.
	 * @param next The next handler for HTTP requests in the pipeline.
	 * @returns An Observable that will not emit anything.
	 */
	public intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
		return next.handle(request).pipe(tap(
			r => {
				if (Object.prototype.hasOwnProperty.call(r, "body") &&
					Object.prototype.hasOwnProperty.call((r as { body: unknown }).body, "alerts")  &&
					(r as {body: {alerts: Array<unknown>}}).body.alerts !== null) { //Ignore alerts with null value) {
					for (const a of (r as { body: { alerts: Array<unknown> } }).body.alerts) {
						this.alertService.newAlert(a as Alert);
					}
				}
			}
		));
	}
}
