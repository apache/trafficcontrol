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
import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { Router, ActivatedRoute, DefaultUrlSerializer } from "@angular/router";

import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { LoggingService } from "../shared/logging.service";
import { AutocompleteValue } from "../utils";

import { ResetPasswordDialogComponent } from "./reset-password-dialog/reset-password-dialog.component";

/**
 * LoginComponent is the controller for the user login form.
 */
@Component({
	selector: "tp-login",
	styleUrls: ["./login.component.scss"],
	templateUrl: "./login.component.html"
})
export class LoginComponent implements OnInit {
	/** The URL to which to redirect users after successful login. */
	private returnURL = "";

	/** Controls if the password is shown in plain text */
	public hide = true;

	/** The password field's autocomplete value. */
	public readonly passwordAutocomplete = AutocompleteValue.CURRENT_PASSWORD;

	/** The user-entered username. */
	public u = "";
	/** The user-entered password. */
	public p: string | null = "";

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly auth: CurrentUserService,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) {
		this.navSvc.headerHidden.next(true);
		this.navSvc.sidebarHidden.next(true);

	}

	/**
	 * Runs initialization, setting up the post-login redirection from the query
	 * string parameters.
	 */
	public async ngOnInit(): Promise<void> {
		const params = this.route.snapshot.queryParamMap;
		this.returnURL = params.get("returnUrl") ?? "core";
		const token = params.get("token");
		if (token) {
			try {
				const response = await this.auth.login(token);
				if (response) {
					this.navSvc.headerHidden.next(false);
					this.navSvc.sidebarHidden.next(false);
					this.router.navigate(["/core/me"], {queryParams: {edit: true, updatePassword: true}});
				}
			} catch (e) {
				this.log.error("token login failed:", e);
			}
		}
	}

	/**
	 * Handles submission of the Login form, and redirects the user back to
	 * their requested page should it be successful. If the user had not yet
	 * requested a page, they will be redirected to `/`
	 */
	public async submitLogin(): Promise<void> {
		if (!this.p) {
			// This shouldn't really be possible, since the value will only be
			// `null` if the control is invalid.
			throw new Error("password is required");
		}
		try {
			const response = await this.auth.login(this.u, this.p);
			if (response) {
				this.navSvc.headerHidden.next(false);
				this.navSvc.sidebarHidden.next(false);
				const tree = new DefaultUrlSerializer().parse(this.returnURL);
				this.router.navigate(tree.root.children.primary.segments.map(s=>s.path), {queryParams: tree.queryParams});
			}
		} catch (err) {
			this.log.error("login failed:", err);
		}
	}

	/** Opens the "reset password" dialog. */
	public resetPassword(): void {
		this.dialog.open(ResetPasswordDialogComponent, {width: "30%"});
	}

}
