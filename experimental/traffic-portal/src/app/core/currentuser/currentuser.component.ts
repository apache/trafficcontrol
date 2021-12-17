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
import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import { faEdit } from "@fortawesome/free-solid-svg-icons";
import { UserService } from "src/app/shared/api";

import { CurrentUser } from "src/app/models";
import {CurrentUserService} from "src/app/shared/currentUser/current-user.service";
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
	public currentUser: CurrentUser | null = null;
	/** Whether or not the page is in 'edit' mode. */
	private editing = false;
	/** Whether or not the page is in 'edit' mode. */
	public get editMode(): boolean {
		return this.editing;
	}
	/** The icon for the 'edit' button. */
	public editIcon = faEdit;
	/**
	 * The editing copy of the current user - used so that you don't need to
	 * reload the page to see accurate information when the edits are cancelled.
	 */
	public editUser: CurrentUser | null = null;

	constructor(
		private readonly auth: CurrentUserService,
		private readonly api: UserService,
		private readonly dialog: MatDialog,
		private readonly route: ActivatedRoute,
		private readonly router: Router
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
					}
				}
			);
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
			console.error("cannot edit null user");
			return;
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
	public submitEdit(e: Event): void {
		e.preventDefault();
		e.stopPropagation();
		if (this.editUser === null) {
			throw new Error("cannot submit edit with null user");
		}

		// There's a separate form for editing passwords, we don't intend to do that here.
		this.editUser.localPasswd = undefined;
		this.editUser.confirmLocalPasswd = undefined;

		this.api.updateCurrentUser(this.editUser).then(
			success => {
				if (success) {
					this.auth.updateCurrentUser().then(
						updated => {
							if (!updated) {
								console.warn("Failed to fetch current user after successful update");
							}
							this.currentUser = this.auth.currentUser;
							this.cancelEdit();
						}
					);
				} else {
					console.warn("Editing the current user failed");
					this.cancelEdit();
				}
			}
		);
	}
}
