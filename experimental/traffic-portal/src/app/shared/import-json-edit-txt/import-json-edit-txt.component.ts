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

import { DatePipe } from "@angular/common";
import { Component, HostListener } from "@angular/core";
import { MatDialogRef } from "@angular/material/dialog";
import { AlertLevel } from "trafficops-types";

import { AlertService } from "../alert/alert.service";

/**
 * Component
 */
@Component({
	selector: "tp-import-json-edit-txt",
	styleUrls: ["./import-json-edit-txt.component.scss"],
	templateUrl: "./import-json-edit-txt.component.html",
})
export class ImportJsonEditTxtComponent {

	/**
	 * Title for the dialog window
	 */
	public title = "Import Profile";

	/**
	 * Allowed import file types
	 */
	public allowedType: string[] = ["application/json", "text/plain"];

	/** Text editor value */
	public inputTxt = "";

	/**  File data imported */
	public fileData = "";

	/** Monitor whether any file is being drag over the dialog */
	public dragOn = false;

	/**
	 * Creates an instance of import json edit txt component.
	 *
	 * @param dialogRef Dialog manager
	 * @param alertService Alert service manager
	 * @param datePipe Default angular date pipe for formating date
	 */
	constructor(
		private readonly dialogRef: MatDialogRef<ImportJsonEditTxtComponent>,
		private readonly alertService: AlertService,
		private readonly datePipe: DatePipe) { }

	/**
	 * Emits the json value for import as profile data
	 */
	public submit(): void {
		this.dialogRef.close(this.inputTxt);
	}

	/**
	 * Hosts listener for drag over
	 *
	 * @param evt Drag events data
	 */
	@HostListener("dragover", ["$event"]) public onDragOver(evt: DragEvent): void {
		evt.preventDefault();
		evt.stopPropagation();

		this.dragOn = true;
	}

	/**
	 * Hosts listener for drag leave
	 *
	 * @param evt Drag events data
	 */
	@HostListener("dragleave", ["$event"]) public onDragLeave(evt: DragEvent): void {
		evt.preventDefault();
		evt.stopPropagation();

		this.dragOn = false;
	}

	/**
	 * Hosts listener for drop
	 *
	 * @param evt Drag events data
	 */
	@HostListener("drop", ["$event"]) public onDrop(evt: DragEvent): void {
		evt.preventDefault();
		evt.stopPropagation();

		this.dragOn = false;

		const file = evt.dataTransfer?.files[0];

		if (!file) {
			return;
		}

		/** Check whether expected file is being uploaded  */
		if (!this.allowedType.find(type => type === file.type)) {
			this.alertService.newAlert({ level: AlertLevel.ERROR, text: "Only JSON, text files are allowed to upload." });
			return;
		}

		/** Format text with data from file data and formated date with date pipe */
		this.fileData = `${file.name} - ${file.size} bytes, last modified: ${this.datePipe.transform(file.lastModified, "MM-dd-yyyy")}`;
		const reader = new FileReader();
		reader.onload = (event): void => {
			this.inputTxt = event.target?.result as string;
		};
		reader.readAsText(file);
	}
}
