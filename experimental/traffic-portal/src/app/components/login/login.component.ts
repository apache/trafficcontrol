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
import { FormControl } from "@angular/forms";
import { Router, ActivatedRoute } from "@angular/router";

import { first } from "rxjs/operators";

import { AuthenticationService } from "../../services";

/**
 * LoginComponent is the controller for the user login form.
 */
@Component({
	selector: "tp-login",
	styleUrls: ["./login.component.scss"],
	templateUrl: "./login.component.html"
})
export class LoginComponent implements OnInit {
	private returnURL: string;

	/** The user-entered username. */
	public u = new FormControl("");
	/** The user-entered password. */
	public p = new FormControl("");

	constructor (
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly auth: AuthenticationService) { }

	/**
	 * Runs initialization, setting up the post-login redirection from the query
	 * string parameters.
	 */
	public ngOnInit(): void {
		this.returnURL = this.route.snapshot.queryParams.returnUrl || "/";
	}

	/**
	 * Handles submission of the Login form, and redirects the user back to their requested page
	 * should it be succesful. If the user had not yet requested a page, they will be redirected to
	 * `/`
	 */
	public submitLogin(): void {
		this.auth.login(this.u.value, this.p.value).pipe(first()).subscribe(
			(response) => {
				if (response) {
					this.router.navigate([this.returnURL]);
				}
			}
		);
	}

}
