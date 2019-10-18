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
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, throwError } from 'rxjs';
import { map, first, catchError } from 'rxjs/operators';

import { Alert } from '../models/alert';

@Injectable({ providedIn: 'root' })
/**
 * This class is responsible for populating an alerts Observable that can be subscribed to by the
 * `AlertComponent`.
*/
export class AlertService {
	public alertsSubject: BehaviorSubject<Alert>;
	public alerts: Observable<Alert>;

	constructor () {
		this.alertsSubject = new BehaviorSubject<Alert>(null);
		this.alerts = this.alertsSubject.asObservable();
	}

	/**
	 * Directly constructs a new UI alert
	 * @param levelOrAlert Either an {@link Alert} or the level of alert
	 * @param text Must be defined if `levelOrAlert` is a String - gives the text of the new alert.
	 * @throws when `levelOrAlert` is a string, but `text` was not provided.
	 */
	public newAlert (levelOrAlert: string | Alert, text?: string) {
		if (typeof levelOrAlert === 'string') {
			if (text === null || text === undefined) {
				throw new Error("Can't pass raw level without raw text!");
			}
			this.alertsSubject.next({level: levelOrAlert, text: text} as Alert);
		} else {
			this.alertsSubject.next(levelOrAlert as Alert);
		}
	}
}
