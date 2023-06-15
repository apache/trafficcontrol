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
import { Component, HostListener, Inject } from "@angular/core";
import { MAT_DIALOG_DATA } from "@angular/material/dialog";
import { AlertLevel } from "trafficops-types";

import { AlertService } from "../alert/alert.service";

/**
 * Contains the structure of the data that ImportJsonTxtComponent accepts.
 */
export interface ImportJsonTxtComponentModel {
	title: string;
}

/**
 * Component for import of json or text
 */
@Component({
	selector: "tp-import-json-txt",
	styleUrls: ["./import-json-txt.component.scss"],
	templateUrl: "./import-json-txt.component.html",
})
export class ImportJsonTxtComponent {

	/**
	 * Allowed import file types
	 */
	public allowedType: string[] = ["application/json", "text/plain"];

	public file: File | null = null;

	/** Text editor value */
	public inputTxt = null;

	/**  File data imported */
	public fileData = "";

	/** Monitor whether any file is being drag over the dialog */
	public dragOn = false;

	public mimeAlertMsg = "Only JSON or text file is allowed.";
	/**
	 * The value of the file input is maintained by extracting drag-and-drop
	 * files and setting the input's value accordingly. Note that when setting
	 * this property, all but the first file are discarded, as it assumes that
	 * multiple selection is not allowed.
	 */
	public get files(): FileList {
		const dt = new DataTransfer();
		if (this.file) {
			dt.items.add(this.file);
		}
		return dt.files;
	}

	public set files(fl: FileList) {
		this.file = fl[0] ?? null;
	}

	/**
	 * Creates an instance of import json edit txt component.
	 *
	 * @param dialogRef Dialog manager
	 * @param alertService Alert service manager
	 * @param datePipe Default angular date pipe for formating date
	 */
	constructor(
		@Inject(MAT_DIALOG_DATA) public readonly data: ImportJsonTxtComponentModel,
		private readonly alertService: AlertService,
		private readonly datePipe: DatePipe) { }

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

		// returns on when there is no file attachment is there
		if (!file) {
			return;
		}

		this.docReader(file);
	}

	/**
	 * Uploads file
	 *
	 * @param event Event object for upload file
	 */
	public uploadFile(event: Event): void {
		const file = (event.target as HTMLInputElement).files?.[0];

		// returns on when there is no file attachment is there
		if (!file) {
			return;
		}
		this.docReader(file);
	  }

	/**
	 * Docs reader
	 *
	 * @param file that is uploaded
	 */
	public docReader(file: File): void {

		/**
		 * Check whether expected file is being uploaded
		 * returns on file wrong file type is uploaded
		 */
		if (!this.allowedType.includes(file.type)) {
			this.alertService.newAlert({ level: AlertLevel.ERROR, text: this.mimeAlertMsg });
			return;
		}

		/** Format text with data from file data and formated date with date pipe */
		this.fileData = `${file.name} - ${file.size} bytes, last modified: ${this.datePipe.transform(file.lastModified, "MM-dd-yyyy")}`;

		const reader = new FileReader();

		reader.onload = (event): void => {
			if(typeof(event.target?.result)==="string"){
				this.inputTxt = JSON.parse(event.target?.result);
			}
		};
		reader.readAsText(file);
	}
}
