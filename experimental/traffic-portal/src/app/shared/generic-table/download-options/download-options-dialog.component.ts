/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { Component, Inject } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { ColDef, CsvExportParams } from "ag-grid-community";

/**
 * Data passed to DownloadOptionsComponent from the grid
 */
export interface DownloadOptionsDialogData {
	name: string;

	columns: ColDef<unknown>[];

	/**
	 * Number of rows selected, should be undefined when only a single.
	 */
	selectedRows: number | undefined;

	visibleRows: number;

	allRows: number;
}

/**
 * Controller for the DownloadOptions component.
 */
@Component({
	selector: "tp-download-options",
	styleUrls: ["./download-options-dialog.component.scss"],
	templateUrl: "./download-options-dialog.component.html"
})
export class DownloadOptionsDialogComponent {
	public fileName: string;

	public visibleColumns: Array<ColDef<unknown>>;

	public columns: Array<ColDef<unknown>>;

	public includeHidden = false;
	public includeHeaders = true;

	public includeFiltered = false;

	public onlySelected = false;

	/**
	 * Number of selected rows, undefined if single selection.
	 */
	public selectedRows: number | undefined;
	public allRows: number;
	public visibleRows: number;

	/** 'C'SV delimiter */
	public seperator = ",";

	constructor(private readonly dialogRef: MatDialogRef<DownloadOptionsDialogComponent,
	CsvExportParams>, @Inject(MAT_DIALOG_DATA) data: DownloadOptionsDialogData) {
		this.fileName = data.name;
		this.selectedRows = data.selectedRows;
		this.allRows = data.allRows;
		this.visibleRows = data.visibleRows;
		this.visibleColumns = [];
		this.columns = [];
		for(const col of data.columns) {
			if(!col.hide) {
				this.visibleColumns.push(col);
			}
			this.columns.push(col);
		}
	}

	/**
	 * Called when submitting the form, converts data into export params.
	 */
	public onSubmit(): void {
		const params: CsvExportParams = {
			allColumns: this.includeHidden,
			columnSeparator: this.seperator,
			exportedRows: this.includeFiltered ? "all" : "filteredAndSorted",
			fileName: `${this.fileName}.csv`,
			onlySelected: this.onlySelected,
			skipColumnHeaders: !this.includeHeaders,
		};
		this.dialogRef.close(params);
	}

}
