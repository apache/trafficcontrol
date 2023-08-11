/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
import { Component, type OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import { ResponseCurrentUser } from "trafficops-types";

import { UserService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import {ThemeManagerService} from "src/app/shared/theme-manager/theme-manager.service";

import { UpdatePasswordDialogComponent } from "./update-password-dialog/update-password-dialog.component";

/**
 * CurrentuserComponent is the controller for the current user's profile page.
 */
@Component({
	selector: "tp-currentuser",
	styleUrls: ["./currentuser.component.scss"],
	templateUrl: "./currentuser.component.html"
})
export class CurrentuserComponent implements OnInit {

	/** The currently logged-in user - or 'null' if not logged-in. */
	public currentUser: ResponseCurrentUser | null = null;
	/** Whether the page is in 'edit' mode. */
	private editing = false;
	/** Whether the page is in 'edit' mode. */
	public get editMode(): boolean {
		return this.editing;
	}
	/**
	 * The editing copy of the current user - used so that you don't need to
	 * reload the page to see accurate information when the edits are cancelled.
	 */
	public editUser: ResponseCurrentUser | null = null;

	constructor(
		private readonly auth: CurrentUserService,
		private readonly api: UserService,
		private readonly dialog: MatDialog,
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly navSvc: NavigationService,
		public readonly themeSvc: ThemeManagerService,
		private readonly log: LoggingService
	) {
		this.currentUser = this.auth.currentUser;
	}

	/**
	 * Runs initialization, setting the currently logged-in user from the
	 * authentication service.
	 */
	public ngOnInit(): void {
		if (this.currentUser === null) {
			this.auth.updateCurrentUser().then(
				r => {
					if (r) {
						this.currentUser = this.auth.currentUser;
						this.navSvc.headerTitle.next(this.currentUser?.username ?? "");
					}
				}
			);
		} else {
			this.navSvc.headerTitle.next(this.currentUser?.username ?? "");
		}
		const edit = this.route.snapshot.queryParamMap.get("edit");
		if (edit === "true") {
			this.edit();
			const updPasswd = this.route.snapshot.queryParamMap.get("updatePassword");
			if (updPasswd === "true") {
				this.updatePassword();
			}
		}
	}

	/**
	 * Handles when the user clicks on the 'edit' button, making the user's
	 * information editable.
	 */
	public edit(): void {
		if (!this.currentUser) {
			throw new Error("cannot edit null user");
		}
		this.editUser = {...this.currentUser};
		this.editing = true;
	}

	/**
	 * Handles when the user click's on the 'cancel' button to cancel edits to
	 * the user's information.
	 */
	public cancelEdit(): void {
		if (!this.currentUser) {
			throw new Error("shouldn't be able to be in edit mode with a null user");
		}
		// It's impossible to be in edit mode with a null user
		this.editUser = {...this.currentUser};
		this.router.navigate(["."], {queryParams: {}, relativeTo: this.route});
		this.editing = false;
	}

	/**
	 * Opens the password change dialog box/form.
	 */
	public updatePassword(): void {
		this.router.navigate(["."], {queryParams: {edit: true, updatePassword: true}, relativeTo: this.route});
		this.dialog.open(UpdatePasswordDialogComponent).afterClosed().subscribe(
			() => {
				this.router.navigate(["."], {queryParams: {edit: true}, relativeTo: this.route});
			}
		);
	}

	/**
	 * Handles submission of the user edit form.
	 *
	 * @param e The form submittal event.
	 */
	public async submitEdit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if (this.editUser === null) {
			throw new Error("cannot submit edit with null user");
		}

		// There's a separate form for editing passwords, we don't intend to do that here.
		const success = await this.api.updateCurrentUser(this.editUser);
		if (success) {
			const updated = await this.auth.updateCurrentUser();
			if (!updated) {
				this.log.warn("Failed to fetch current user after successful update");
			}
			this.currentUser = this.auth.currentUser;
			this.cancelEdit();
		} else {
			this.log.warn("Editing the current user failed");
			this.cancelEdit();
		}
	}

	/**
	 * Checks if the form's user has a "bottom-level" address, meaning any
	 * combination of state/province, postal code, city, and/or country.
	 *
	 * @returns `true` if the user has a "bottom-level" address, `false`
	 * otherwise.
	 */
	public hasBottomAddress(): boolean {
		if (!this.currentUser) {
			return false;
		}
		const {city, country, stateOrProvince, postalCode} = this.currentUser;
		return !!(city || country || stateOrProvince || postalCode);
	}
}
