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
import type { MatSelectChange } from "@angular/material/select";
import { ActivatedRoute, Router } from "@angular/router";
import type { PostRequestUser, ResponseRole, ResponseTenant, ResponseUser, User } from "trafficops-types";

import { UserService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * UserDetailsComponent is the controller for the page for viewing/editing a
 * user.
 */
@Component({
	selector: "tp-user-details",
	styleUrls: ["./user-details.component.scss"],
	templateUrl: "./user-details.component.html"
})
export class UserDetailsComponent implements OnInit {

	public user!: ResponseUser | PostRequestUser;
	public roles = new Array<ResponseRole>();
	public tenants = new Array<ResponseTenant>();
	public new = false;

	constructor(
		private readonly userService: UserService,
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly currentUserService: CurrentUserService,
		private readonly log: LoggingService,
		private readonly navSvc: NavigationService,
	) { }

	/** Angular lifecycle hook */
	public async ngOnInit(): Promise<void> {
		const rolesAndTenants = Promise.all([
			this.userService.getRoles().then(rs=>this.roles=rs),
			this.userService.getTenants().then(ts=>this.tenants=ts)
		]);
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			this.log.error("missing required route parameter 'id'");
			return;
		}
		await rolesAndTenants;
		this.new = ID === "new";

		if (this.new) {
			this.setTitle();
			this.new = true;
			this.user = {
				confirmLocalPasswd: "",
				email: "user@example.com",
				fullName: "",
				localPasswd: "",
				role: this.currentUserService.currentUser?.role ?? "",
				tenantId: this.currentUserService.currentUser?.tenantId ?? 1,
				username: "",
			};
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			this.log.error("route parameter 'id' was non-number:", ID);
			return;
		}
		this.user = await this.userService.getUsers(numID);
		this.setTitle();
	}

	/**
	 * Used to tell whether the form is for adding or editing a user.
	 *
	 * @param _ The user represented by the component. This is totally
	 * unnecessary for calculating the result, just needed to make the compiler
	 * happy.
	 * @returns Whether the form is for a new user (`true`) or an existing user
	 * (`false`).
	 */
	public isNew(_?: User): _ is PostRequestUser {
		return this.new;
	}

	/**
	 * Sets the headerTitle based on current User state.
	 *
	 * @private
	 */
	private setTitle(): void {
		const title = this.new ? "New User" : `User: ${this.user.username}`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * Handler for the user edit form submission.
	 *
	 * @param e The form submission event. Its default behavior of sending an
	 * HTTP request is disabled.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if (this.isNew(this.user)) {
			this.user = await this.userService.createUser(this.user);
			this.new = false;
			await this.router.navigate(["core/users", this.user.id]);
		}
		this.user = await this.userService.updateUser(this.user);
		this.setTitle();
	}

	/**
	 * Gets the Role of the User the form is editing.
	 *
	 * @returns The user's current Role.
	 */
	public role(): ResponseRole | null {
		if (this.isNew(this.user)) {
			return null;
		}
		const role = this.roles.find(r=>r.name === this.user.role);
		if (!role) {
			throw new Error(`user's Role "${this.user.role}" does not exist`);
		}
		return role;
	}

	/**
	 * Gets the Tenant of the User the form is editing.
	 *
	 * @returns The user's current Tenant.
	 */
	public tenant(): ResponseTenant | null {
		if (this.isNew(this.user)) {
			return null;
		}
		const tenant = this.tenants.find(t=>t.id === this.user.tenantId);
		if (!tenant) {
			throw new Error(`user's Tenant "${this.user.tenant}" (#${this.user.tenantId}) does not exist`);
		}
		return tenant;
	}

	/**
	 * Handles changes to the role selection by updating the `role` and
	 * `rolename` properties of the form's User accordingly.
	 *
	 * @param r The Role selected by the user.
	 */
	public updateRole(r: MatSelectChange & {value: ResponseRole}): void {
		this.user.role = r.value.name;
	}

	/**
	 * Handles changes to the tenant selection by updating the `tenant` and
	 * `tenantId` properties of the form's User accordingly.
	 *
	 * @param t The Tenant selected by the user.
	 */
	public updateTenant(t: MatSelectChange & {value: ResponseTenant}): void {
		this.user.tenantId = t.value.id;
		if (!this.isNew(this.user)) {
			this.user.tenant = t.value.name;
		}
	}
}
