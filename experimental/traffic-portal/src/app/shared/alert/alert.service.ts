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
import { Injectable } from "@angular/core";
import { BehaviorSubject, type Observable } from "rxjs";
import type { Alert, AlertLevel } from "trafficops-types";

/**
 * This class is responsible for populating an alerts Observable that can be
 * subscribed to by the `AlertComponent`.
 */
@Injectable({
	providedIn: "root"
})
export class AlertService {
	/** A BehaviorSubject that emits Alerts. */
	private readonly alertsSubject: BehaviorSubject<Alert | null>;
	/** An Observable that emits Alerts. */
	public alerts: Observable<Alert | null>;

	/**
	 * Constructor.
	 */
	constructor() {
		this.alertsSubject = new BehaviorSubject<Alert | null>(null);
		this.alerts = this.alertsSubject.asObservable();
	}

	/**
	 * Directly constructs a new UI alert.
	 *
	 * @param level The level of the Alert.
	 * @param text The message content of the Alert.
	 */
	public newAlert(level: AlertLevel, text: string): void;
	/**
	 * Directly constructs a new UI alert.
	 *
	 * @param alert The Alert to be raised.
	 */
	public newAlert(alert: Alert): void;
	/**
	 * Directly constructs a new UI alert
	 *
	 * @param levelOrAlert Either an {@link Alert} or the level of alert
	 * @param text Must be defined if `levelOrAlert` is a String - gives the text of the new alert.
	 * @throws when `levelOrAlert` is a string, but `text` was not provided.
	 */
	public newAlert(levelOrAlert: AlertLevel | Alert, text?: string): void {
		if (typeof levelOrAlert === "string") {
			if (text === undefined) {
				throw new Error("Can't pass raw level without raw text!");
			}
			this.alertsSubject.next({level: levelOrAlert, text});
		} else {
			this.alertsSubject.next(levelOrAlert);
		}
	}
}
