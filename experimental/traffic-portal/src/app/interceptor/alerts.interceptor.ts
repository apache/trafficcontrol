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

import { Observable } from "rxjs";
import { tap } from "rxjs/operators";

import { Alert } from "../models/alert";
import { AlertService } from "../services";

/**
 * This class intercepts any and all alerts contained in API responses and
 * passes them to the `AlertService` for display to the user.
 */
@Injectable()
export class AlertInterceptor implements HttpInterceptor {
	constructor (private readonly alertService: AlertService) {}

	/**
	 * Intercepts HTTP responses and checks for any alerts.
	 */
	public intercept (request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
		return next.handle(request).pipe(tap(
			r => {
				/* tslint:disable */
				if (r.hasOwnProperty('body') && r['body'].hasOwnProperty('alerts')) {
					for (const a of r['body']['alerts']) {
						/* tslint:enable */
						this.alertService.alertsSubject.next(a as Alert);
					}
				}
			}
		));
	}
}
