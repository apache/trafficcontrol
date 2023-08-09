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
import { LoggingService } from "../logging.service";

/**
 * Contains the structure of the data that {@link ImportJsonTxtComponent}
 * accepts.
 */
export interface ImportJsonTxtComponentModel {
	title: string;
}

/**
 * Component for import of JSON or text files.
 */
@Component({
	selector: "tp-import-json-txt",
	styleUrls: ["./import-json-txt.component.scss"],
	templateUrl: "./import-json-txt.component.html",
})
export class ImportJsonTxtComponent {

	/**
	 * Allowed import file types.
	 */
	public readonly allowedType: readonly string[] = ["application/json", "text/plain"];

	public file: File | null = null;

	/** Text editor value */
	public inputTxt = "";

	/**  File data imported */
	public fileData = "";

	/** Monitor whether any file is being drag over the dialog */
	public dragOn = false;

	public readonly mimeAlertMsg = "Only JSON or text files are allowed.";

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
	 * Constructor.
	 *
	 * @param data Data passed as input to the component.
	 * @param dialogRef Angular dialog service.
	 * @param alertService Alerts service.
	 * @param datePipe Default Angular pipe used for formatting dates.
	 * @param log Logging service.
	 */
	constructor(
		@Inject(MAT_DIALOG_DATA) public readonly data: ImportJsonTxtComponentModel,
		private readonly alertService: AlertService,
		private readonly datePipe: DatePipe,
		private readonly log: LoggingService,
	) { }

	/**
	 * Hosts listener for drag over
	 *
	 * @param evt Drag events data
	 */
	@HostListener("dragover", ["$event"])
	public onDragOver(evt: DragEvent): void {
		evt.preventDefault();
		evt.stopPropagation();

		this.dragOn = true;
	}

	/**
	 * Hosts listener for drag leave
	 *
	 * @param evt Drag events data
	 */
	@HostListener("dragleave", ["$event"])
	public onDragLeave(evt: DragEvent): void {
		evt.preventDefault();
		evt.stopPropagation();

		this.dragOn = false;
	}

	/**
	 * Hosts listener for drop
	 *
	 * @param evt Drag events data
	 */
	@HostListener("drop", ["$event"])
	public onDrop(evt: DragEvent): void {
		evt.preventDefault();
		evt.stopPropagation();

		this.dragOn = false;
		if (!evt.dataTransfer) {
			return;
		}

		this.files = evt.dataTransfer.files;
		this.docReader();
	}

	/**
	 * Uploads file
	 *
	 * @param event Event object for upload file
	 */
	public uploadFile(event: Event): void {
		if (!(event.target instanceof HTMLInputElement) || !event.target.files) {
			this.log.warn("file uploading triggered on non-file-input element:", event.target);
			return;
		}

		this.files = event.target.files;
		this.docReader();
	}

	/**
	 * Docs reader
	 *
	 * @param file that is uploaded
	 */
	private docReader(): void {
		if (!this.file) {
			return;
		}

		/**
		 * Check whether expected file is being uploaded
		 * returns on file wrong file type is uploaded
		 */
		if (!this.allowedType.includes(this.file.type)) {
			this.alertService.newAlert({ level: AlertLevel.ERROR, text: this.mimeAlertMsg });
			return;
		}

		/** Format text with data from file data and formated date with date pipe */
		const dateStr = this.datePipe.transform(this.file.lastModified, "MM-dd-yyyy");
		this.fileData = `${this.file.name} - ${this.file.size} bytes, last modified: ${dateStr}`;

		const reader = new FileReader();
		reader.addEventListener("load",
			event => {
				if(typeof(event.target?.result)==="string"){
					this.inputTxt = JSON.parse(event.target.result);
				}
			}
		);
		reader.readAsText(this.file);
	}
}
