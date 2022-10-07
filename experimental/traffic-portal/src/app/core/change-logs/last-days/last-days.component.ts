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
import { Component, Inject} from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";

/**
 * LastDaysComponent contains logic to change the number of days to view change logs.
 */
@Component({
	selector: "tp-last-days",
	styleUrls: ["./last-days.component.scss"],
	templateUrl: "./last-days.component.html"
})
export class LastDaysComponent {
	public days: string;

	constructor(private readonly dialogRef: MatDialogRef<LastDaysComponent>, @Inject(MAT_DIALOG_DATA) private readonly lastDays: string) {
		this.days = this.lastDays;
	}

	/**
	 * Emits from the new number of days.
	 */
	public submit(): void {
		this.dialogRef.close(this.days);
	}

	/**
	 * Emits from the `done` Output indicating the action was cancelled.
	 */
	public cancel(): void {
		this.dialogRef.close(this.lastDays);
	}
}
