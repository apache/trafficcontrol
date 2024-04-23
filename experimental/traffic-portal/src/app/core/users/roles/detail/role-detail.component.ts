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
import { ActivatedRoute, Router } from "@angular/router";
import { ResponseRole } from "trafficops-types";

import { UserService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * RoleDetailComponent is the controller for the Role add/edit form.
 */
@Component({
	selector: "tp-role-detail",
	templateUrl: "./role-detail.component.html"
})
export class RoleDetailComponent implements OnInit {
	public new = false;
	public permissions = "";
	public role!: ResponseRole;

	/**
	 * This caches the original name of the Role, so that updates can be
	 * made.
	 */
	private name = "";

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly userService: UserService,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) { }

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		const role = this.route.snapshot.paramMap.get("name");
		if (role === null) {
			this.setHeader("New Role");
			this.new = true;
			this.role = {
				description: "",
				name: "",
				permissions: []
			};
			return;
		}

		this.role = await this.userService.getRoles(role);
		this.name = this.role.name;
		this.permissions = this.role.permissions?.join("\n")??"";
		this.navSvc.headerTitle.next(`Role: ${this.name}`);
	}

	/**
	 * Sets the value of the header text, and caches the Role's initial
	 * name.
	 *
	 * @param name The name of the current Role (before editing).
	 */
	private setHeader(name: string): void {
		this.name = name;
		this.navSvc.headerTitle.next(`Role: ${name}`);
	}

	/**
	 * Deletes the current Role.
	 */
	public async deleteRole(): Promise<void> {
		if (this.new) {
			this.log.error("Unable to delete new role");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete role ${this.role.name}`,
				title: "Confirm Delete"}
		});
		const result = await ref.afterClosed().toPromise();
		if(result) {
			await this.userService.deleteRole(this.role);
			await this.router.navigate(["core/roles"]);
		}
	}

	/**
	 * Updates permissions list from a string to an array.
	 */
	public async updatePermissions(): Promise<void> {
		this.role.permissions = this.permissions.split("\n");
	}

	/**
	 * Submits new/updated Role.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		await this.updatePermissions();
		if(this.new) {
			this.role = await this.userService.createRole(this.role);
			this.new = false;
		} else {
			this.role = await this.userService.updateRole(this.name, this.role);
			this.navSvc.headerTitle.next(`Role: ${this.role.name}`);
		}
		await this.router.navigate(["core/roles", this.role.name], {replaceUrl: true});
	}

}
