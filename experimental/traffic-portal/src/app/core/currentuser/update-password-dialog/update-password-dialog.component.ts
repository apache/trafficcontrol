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
import { Component } from "@angular/core";
import { MatDialogRef } from "@angular/material/dialog";
import { Subject } from "rxjs";

import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { LoggingService } from "src/app/shared/logging.service";

/**
 * This is the controller for the "Update Password" dialog box/form.
 */
@Component({
	selector: "tp-update-password-dialog",
	styleUrls: ["./update-password-dialog.component.scss"],
	templateUrl: "./update-password-dialog.component.html",
})
export class UpdatePasswordDialogComponent {
	/** The new password. */
	public password = "";
	/** The new password repeated for confirmation. */
	public confirm = "";

	/** A Subscribable that tracks the validity ofthe password confirmation field. */
	public readonly confirmValid = new Subject<string>();

	constructor(
		private readonly dialog: MatDialogRef<UpdatePasswordDialogComponent>,
		private readonly auth: CurrentUserService,
		private readonly log: LoggingService,
	) {}

	/**
	 * Cancels the password update, closing the dialog box.
	 */
	public cancel(): void {
		this.dialog.close();
	}

	/**
	 * Handles submission of the form, checking that the passwords match before
	 * sending them to the server.
	 *
	 * @param event The form submission event, which must have its default
	 * prevented.
	 */
	public async submit(event: Event): Promise<void> {
		event.preventDefault();
		event.stopPropagation();

		if (this.confirm !== this.password) {
			this.confirmValid.next("Passwords do not match");
			return;
		}

		if (!this.auth.currentUser) {
			this.log.error("Cannot update null user");
			return;
		}

		const user = {
			...this.auth.currentUser,
			confirmLocalPasswd: this.confirm,
			localPasswd: this.password,
		};

		user.localPasswd = this.password;
		user.confirmLocalPasswd = this.confirm;
		return this.auth.saveCurrentUser(user).then(
			success => {
				if (success) {
					this.dialog.close(true);
				}
			}
		);
	}

}
