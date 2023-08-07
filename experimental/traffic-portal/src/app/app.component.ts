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
import { Component, Inject, OnInit, Optional, PLATFORM_ID, TransferState, makeStateKey } from "@angular/core";
import { Router } from "@angular/router";
import { ResponseCurrentUser } from "trafficops-types";

import { CurrentUserService } from "src/app/shared/current-user/current-user.service";

export const LOCAL_TPV1_URL = "tp_v1_url";

/**
 * The most basic component that contains everything else. This should be kept pretty simple.
 */
@Component({
	selector: "app-root",
	styleUrls: ["./app.component.scss"],
	templateUrl: "./app.component.html",
})
export class AppComponent implements OnInit {

	/** The currently logged-in user */
	public currentUser: ResponseCurrentUser | null = null;

	constructor(private readonly router: Router, private readonly auth: CurrentUserService,
		@Inject(PLATFORM_ID) private readonly platformId: object,
		@Optional() @Inject("TP_V1_URL") public tpv1url: string,
		private readonly transferState: TransferState) {
		const storeKey = makeStateKey<string>("messageKey");

		// get data from transferState if browser side
		if (isPlatformBrowser(this.platformId)) {
			this.tpv1url = this.transferState.get(storeKey, "https://localhost");
			window.localStorage.setItem(LOCAL_TPV1_URL, this.tpv1url);
		} else { // server side: get provided tpv1 url and store in in transfer state
			this.transferState.set(storeKey, this.tpv1url);
		}
	}

	/**
	 * Logs the currently logged-in user out.
	 */
	public async logout(): Promise<void> {
		this.auth.logout();
		await this.router.navigate(["/login"]);
	}

	/**
	 * Sets up the current user.
	 */
	public ngOnInit(): void {
		this.auth.userChanged.subscribe(user => {
			this.currentUser = user;
		});
	}

}
