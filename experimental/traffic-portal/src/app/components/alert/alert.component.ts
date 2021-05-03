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
import { MatSnackBar } from "@angular/material/snack-bar";
import { Component, OnDestroy } from "@angular/core";
import { Subscription } from "rxjs";

import { AlertService } from "../../services";

/**
 * AlertComponent is the controller for alert message popups.
 */
@Component({
	selector: "tp-alert",
	styleUrls: ["./alert.component.scss"],
	templateUrl: "./alert.component.html"
})
export class AlertComponent implements OnDestroy {

	/** Internal subscription to the AlertService's alerts observable. */
	private readonly subscription: Subscription;

	/** The duration for which Alerts will linger until dismissed. `undefined` means forever. */
	public duration: number | undefined = 10000;

	/**
	 * Constructor.
	 */
	constructor(
		private readonly alerts: AlertService,
		private readonly snackBar: MatSnackBar
	) {
		this.subscription = this.alerts.alerts.subscribe(
			a => {
				if (a) {
					if (a.text === "") {
						a.text = "Unknown";
					}
					switch (a.level) {
						case "success":
							console.log("alert: ", a.text);
							break;
						case "info":
							console.log("alert: ", a.text);
							break;
						case "warning":
							console.warn("alert: ", a.text);
							break;
						case "error":
							console.error("alert: ", a.text);
							break;
						default:
							console.log("unknown alert: ", a.text);
							break;
					}
					this.snackBar.open(a.text, "dismiss", {duration: this.duration, verticalPosition: "top"});
				}
			},
			e => {
				console.error("Error in alerts subscription:", e);
			}
		);
	}

	/**
	 * Clears away the currently visible Alert.
	 */
	public clear(): void {
		this.snackBar.dismiss();
	}

	/**
	 * Cleans up persistent resources in the component.
	 */
	public ngOnDestroy(): void {
		this.subscription.unsubscribe();
	}
}
