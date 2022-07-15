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
import { formatDate } from "@angular/common";
import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { ColDef } from "ag-grid-community";
import { BehaviorSubject } from "rxjs";

import { ChangeLogsService } from "src/app/api/change-logs.service";
import { LastDaysComponent } from "src/app/core/change-logs/last-days/last-days.component";
import { ChangeLog } from "src/app/models/change-logs";
import { TableTitleButton } from "src/app/shared/generic-table/generic-table.component";
import { TpHeaderService } from "src/app/shared/tp-header/tp-header.service";

/**
 * AugmentedChangeLog has fields for access to processed times.
 */
interface AugmentedChangeLog extends ChangeLog {
	longTime: string;
	relativeTime: string;
}

/**
 * Converts a changelog to an augmented one.
 *
 * @param cl Changelog to convert.
 * @returns Converted changelog
 */
export function augment(cl: ChangeLog): AugmentedChangeLog {
	const aug = {longTime: "", relativeTime: "", ...cl};

	const d = new Date(aug.lastUpdated);
	const delta = new Date().getTime() - d.getTime();
	const SEC = 1000;
	const MIN = SEC * 60;
	const HOUR = MIN * 60;
	const DAY = HOUR * 24;
	const WEEK = DAY * 7;
	const MONTH = WEEK * 4;
	const YEAR = MONTH * 12;
	if (delta > YEAR) {
		aug.relativeTime = `${(delta / YEAR).toFixed(2)} years ago`;
	} else if (delta > MONTH) {
		aug.relativeTime =  `${(delta / MONTH).toFixed(2)} months ago`;
	} else if (delta > WEEK) {
		aug.relativeTime = `${(delta/ WEEK).toFixed(2)} weeks ago`;
	} else if (delta > DAY) {
		aug.relativeTime = `${(delta / DAY).toFixed(2)} days ago`;
	} else if (delta > HOUR) {
		aug.relativeTime = `${(delta / HOUR).toFixed(2)} hours ago`;
	} else if (delta > MIN) {
		aug.relativeTime = `${(delta / MIN).toFixed(2)} minutes ago`;
	} else {
		aug.relativeTime = `${(delta / SEC).toFixed(0)} seconds ago`;
	}

	aug.longTime = formatDate(aug.lastUpdated, "long", "en-US");
	return aug;
}

/**
 *  ChangeLogsComponent is the controller for the change logs page
 */
@Component({
	selector: "tp-change-logs",
	styleUrls: ["./change-logs.component.scss"],
	templateUrl: "./change-logs.component.html"
})
export class ChangeLogsComponent implements OnInit {
	/** Emits changes to the fuzzy search text. */
	public fuzzySubj: BehaviorSubject<string>;

	/** The current search text. */
	public searchText = "";

	public changeLogs: Promise<Array<AugmentedChangeLog>> | null = null;

	public lastDays = 7;

	public titleBtns: Array<TableTitleButton> = [
		{
			action: "lastDays",
			text: `Last ${this.lastDays} days`,
		}
	];

	public columnDefs: Array<ColDef> = [
		{
			field: "relativeTime",
			headerName: "Occurred",
		},
		{
			field: "longTime",
			headerName: "Created (UTC)",
		},
		{
			field: "user",
			headerName: "User"
		},
		{
			field: "level",
			headerName: "Level",
			hide: true
		},
		{
			field: "message",
			headerName: "Message"
		}
	];

	/** Whether user data is still loading. */
	public loading = true;

	constructor(private readonly headerSvc: TpHeaderService, private readonly api: ChangeLogsService,
		private readonly route: ActivatedRoute, private readonly dialog: MatDialog) {
		this.fuzzySubj = new BehaviorSubject<string>("");
	}

	/**
	 * Loads the table data based on lastDays.
	 */
	public async loadData(): Promise<void> {
		this.loading = true;
		this.changeLogs = this.api.getChangeLogs({days: this.lastDays.toString()}).then(data => {
			this.loading = false;
			return data.map(augment);
		});
	}

	/**
	 * Angular lifecycle hook
	 */
	public async ngOnInit(): Promise<void> {
		this.headerSvc.headerTitle.next("Change Logs");
		await this.loadData();

		this.route.queryParamMap.subscribe(
			m => {
				const search = m.get("search");
				if (search) {
					this.fuzzySubj.next(search);
				}
				this.searchText = search ?? "";
			}
		);
	}

	/**
	 * handles when a title button is event is emitted
	 *
	 * @param action which button was pressed
	 */
	public async handleTitleButton(action: string): Promise<void> {
		switch(action){
			case "lastDays":
				const ref = this.dialog.open(LastDaysComponent, {
					data: this.lastDays
				});
				ref.afterClosed().subscribe(result => {
					if (result) {
						this.lastDays = +result.toString();
						this.loadData();
						this.titleBtns = [
							{
								action: "lastDays",
								text: `Last ${this.lastDays} days`,
							}
						];
					}
				});
				break;
		}
	}

	/**
	 * Updates the "search" query parameter in the URL every time the search
	 * text input changes.
	 */
	public updateURL(): void {
		this.fuzzySubj.next(this.searchText);
	}
}
