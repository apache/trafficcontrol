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

import { Component, type OnInit } from "@angular/core";
import type { MatSelectChange } from "@angular/material/select";
import { ActivatedRoute } from "@angular/router";

import { UserService } from "src/app/api";
import type { Role, Tenant, User } from "src/app/models";

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

	public user!: User;
	public roles = new Array<Role>();
	public tenants = new Array<Tenant>();

	constructor(private readonly userService: UserService, private readonly route: ActivatedRoute) {
	}

	/** Angular lifecycle hook */
	public async ngOnInit(): Promise<void> {
		const rolesAndTenants = Promise.all([
			this.userService.getRoles().then(rs=>this.roles=rs),
			this.userService.getTenants().then(ts=>this.tenants=ts)
		]);
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			console.error("missing required route parameter 'id'");
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			console.error("route parameter 'id' was non-number:", ID);
			return;
		}
		await rolesAndTenants;
		this.user = await this.userService.getUsers(numID);
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
		this.user = await this.userService.updateUser(this.user);
	}

	/**
	 * Gets the Role of the User the form is editing.
	 *
	 * @returns The user's current Role.
	 */
	public role(): Role {
		const role = this.roles.find(r=>r.id === this.user.role);
		if (!role) {
			throw new Error(`user's Role "${this.user.rolename}" (#${this.user.role}) does not exist`);
		}
		return role;
	}

	/**
	 * Gets the Tenant of the User the form is editing.
	 *
	 * @returns The user's current Tenant.
	 */
	public tenant(): Tenant {
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
	public updateRole(r: MatSelectChange & {value: Role}): void {
		this.user.role = r.value.id;
		this.user.rolename = r.value.name;
	}

	/**
	 * Handles changes to the tenant selection by updating the `tenant` and
	 * `tenantId` properties of the form's User accordingly.
	 *
	 * @param t The Tenant selected by the user.
	 */
	public updateTenant(t: MatSelectChange & {value: Tenant}): void {
		this.user.tenantId = t.value.id;
		this.user.tenant = t.value.name;
	}
}
