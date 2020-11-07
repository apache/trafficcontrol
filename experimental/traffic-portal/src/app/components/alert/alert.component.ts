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
import { isPlatformBrowser } from "@angular/common";
import { Component, Inject, OnInit, PLATFORM_ID } from "@angular/core";

import { Alert } from "../../models";
import { AlertService } from "../../services";

/**
 * AlertComponent is the controller for alert message popups.
 */
@Component({
	selector: "tp-alert",
	styleUrls: ["./alert.component.scss"],
	templateUrl: "./alert.component.html"
})
export class AlertComponent implements OnInit {

	/** The raw HTML element that serves as the base of the component. */
	public dialogElement: HTMLDialogElement;

	/** The Alert being displayed. */
	public alert: Alert | null;

	constructor (private readonly alerts: AlertService, @Inject(PLATFORM_ID) private readonly PLATFORM: Object) { }

	/**
	 * Runs initialization code, setting up the current alert from the alerts
	 * service.
	 */
	public ngOnInit(): void {
		if (!isPlatformBrowser(this.PLATFORM)) {
			return;
		}
		this.dialogElement = document.getElementById("alert") as HTMLDialogElement;
		this.alerts.alerts.subscribe(
			(a: Alert) => {
				if (a) {
					this.alert = a;
					if (a.text === "") {
						a.text = "Unknown";
					}
					switch (a.level) {
						case "success":
							console.log("alert: ", a.text);
							break;
						case "info":
							console.debug("alert: ", a.text);
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
					this.dialogElement.showModal();
				}
			}
		);
	}

	/**
	 * Closes the dialog.
	 */
	public close(): void {
		this.dialogElement.close();
		this.alert = null;
	}
}
