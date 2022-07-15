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
