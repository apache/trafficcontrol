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
import { MatDialogRef } from "@angular/material/dialog";

import { UserService } from "src/app/api";
import { Role, Tenant } from "src/app/models";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";

/**
 * Controller for a dialog that opens to register a new user.
 */
@Component({
	selector: "tp-user-registration-dialog",
	styleUrls: ["./user-registration-dialog.component.scss"],
	templateUrl: "./user-registration-dialog.component.html"
})
export class UserRegistrationDialogComponent implements OnInit {

	public roles = new Array<Role>();
	public tenants = new Array<Tenant>();

	public role!: Role;
	public tenant!: Tenant;
	public email = "";

	constructor(
		private readonly userService: UserService,
		private readonly auth: CurrentUserService,
		private readonly dialogRef: MatDialogRef<UserRegistrationDialogComponent>
	) { }

	/**
	 * Sets up Role and Tenant data using the API.
	 */
	public ngOnInit(): void {
		this.userService.getRoles().then(
			rs => {
				this.roles = rs;
				for (const role of rs) {
					if (role.id === this.auth.currentUser?.role) {
						this.role = role;
					}
				}
			}
		);
		this.userService.getTenants().then(
			ts => {
				this.tenants = ts;
				for (const tenant of ts) {
					if (tenant.id === this.auth.currentUser?.tenantId) {
						this.tenant = tenant;
					}
				}
			}
		);
	}

	/**
	 * Submits the API request to create the user.
	 *
	 * @param e The form submittal event that triggered calling this method. Its
	 * default is prevented, and its propagation stopped.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();

		try {
			await this.userService.registerUser(this.email, this.role, this.tenant);
			this.dialogRef.close();
		} catch (err) {
			console.error("failed to register user:", err);
		}
	}
}
