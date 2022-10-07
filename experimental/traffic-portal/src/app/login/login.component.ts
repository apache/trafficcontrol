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
import { Router, ActivatedRoute } from "@angular/router";

import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import {TpHeaderService} from "src/app/shared/tp-header/tp-header.service";

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
		private readonly headerSvc: TpHeaderService
	) { }

	/**
	 * Runs initialization, setting up the post-login redirection from the query
	 * string parameters.
	 */
	public ngOnInit(): void {
		this.headerSvc.headerHidden.next(true);
		const params = this.route.snapshot.queryParamMap;
		this.returnURL = params.get("returnUrl") ?? "core";
		const token = params.get("token");
		if (token) {
			this.auth.login(token).then(
				response => {
					if (response) {
						this.headerSvc.headerHidden.next(false);
						this.router.navigate(["/core/me"], {queryParams: {edit: true, updatePassword: true}});
					}
				},
				err => {
					console.error("login with token failed:", err);
				}
			);
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
				this.headerSvc.headerHidden.next(false);
				this.router.navigate([this.returnURL]);
			}
		} catch (err) {
			console.error("login failed:", err);
		}
	}

	/** Opens the "reset password" dialog. */
	public resetPassword(): void {
		this.dialog.open(ResetPasswordDialogComponent, {width: "30%"});
	}

}
