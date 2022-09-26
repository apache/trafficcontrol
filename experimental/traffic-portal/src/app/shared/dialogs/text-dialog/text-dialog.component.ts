import { Component, Inject } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";

/**
 * Contains the structure of the data the TextDialogComponent expects
 */
export interface TextDialogData {
	title: string;
	message: string;
}

/**
 * TextDialogComponent contains code for displaying a simple text mat-dialog.
 */
@Component({
	selector: "tp-text-dialog",
	styleUrls: ["./text-dialog.component.scss"],
	templateUrl: "./text-dialog.component.html"
})
export class TextDialogComponent {

	constructor(private readonly dialogRef: MatDialogRef<TextDialogComponent>,
		@Inject(MAT_DIALOG_DATA) public readonly dialogData: TextDialogData) {
	}

	/**
	 * Closes the dialog
	 */
	public close(): void {
		this.dialogRef.close();
	}
}
