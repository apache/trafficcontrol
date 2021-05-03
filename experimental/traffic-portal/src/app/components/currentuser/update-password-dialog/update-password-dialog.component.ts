import { Component } from "@angular/core";
import { MatDialogRef } from "@angular/material/dialog";
import { Subject } from "rxjs";

import { AuthenticationService } from "src/app/services";
import { UserService } from "src/app/services/api";

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
		private readonly auth: AuthenticationService,
		private readonly api: UserService
	) { }

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

		const user = this.auth.currentUser;
		if (!user) {
			console.error("Cannot update null user");
			return;
		}

		user.localPasswd = this.password;
		user.confirmLocalPasswd = this.confirm;
		return this.api.updateCurrentUser(user).then(
			success => {
				if (success) {
					this.dialog.close(true);
				}
			}
		);
	}

}
