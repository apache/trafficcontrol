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

import { UserService } from "src/app/api";

/**
 * The controller for the "reset password" dialog.
 */
@Component({
	selector: "tp-reset-password-dialog",
	styleUrls: ["./reset-password-dialog.component.scss"],
	templateUrl: "./reset-password-dialog.component.html"
})
export class ResetPasswordDialogComponent {

	/** The email address used for password recovery. */
	public email = "";

	constructor(private readonly api: UserService, private readonly dialogRef: MatDialogRef<ResetPasswordDialogComponent>) { }

	/**
	 * Handles submission of the dialog's form.
	 *
	 * @param event The form submission event.
	 */
	public submit(event: Event): void {
		event.preventDefault();
		event.stopPropagation();
		this.api.resetPassword(this.email);
		this.dialogRef.close();
	}

}
