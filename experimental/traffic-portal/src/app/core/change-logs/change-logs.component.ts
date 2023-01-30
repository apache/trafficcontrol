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
import { Log } from "trafficops-types";

import { ChangeLogsService } from "src/app/api/change-logs.service";
import { LastDaysComponent } from "src/app/core/change-logs/last-days/last-days.component";
import { TextDialogComponent } from "src/app/shared/dialogs/text-dialog/text-dialog.component";
import {
	ContextMenuActionEvent,
	ContextMenuItem,
	TableTitleButton
} from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { relativeTimeString } from "src/app/utils";

/**
 * AugmentedChangeLog has fields for access to processed times.
 */
interface AugmentedChangeLog extends Log {
	longTime: string;
	relativeTime: string;
}

/**
 * Converts a changelog to an augmented one.
 *
 * @param cl Changelog to convert.
 * @returns Converted changelog
 */
export function augment(cl: Log): AugmentedChangeLog {
	const aug = {longTime: "", relativeTime: "", ...cl};

	aug.relativeTime = relativeTimeString(new Date().getTime() - new Date(aug.lastUpdated).getTime());
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

	public contextMenuOptions: Array<ContextMenuItem<AugmentedChangeLog>> = [
		{
			action: "viewChangeLog",
			name: "Expand Log"
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

	constructor(private readonly navSvc: NavigationService, private readonly api: ChangeLogsService,
		private readonly route: ActivatedRoute, private readonly dialog: MatDialog) {
		this.fuzzySubj = new BehaviorSubject<string>("");
		this.navSvc.headerTitle.next("Change Logs");
	}

	/**
	 * Loads the table data based on lastDays.
	 */
	public async loadData(): Promise<void> {
		this.loading = true;
		const data = await this.api.getChangeLogs({days: this.lastDays.toString()});
		this.changeLogs = Promise.resolve(data.map(augment));
		this.loading = false;
	}

	/**
	 * Angular lifecycle hook
	 */
	public async ngOnInit(): Promise<void> {
		this.route.queryParamMap.subscribe(
			m => {
				const search = m.get("search");
				if (search) {
					this.fuzzySubj.next(search);
				}
				this.searchText = search ?? "";
			}
		);

		await this.loadData();
	}

	/**
	 * handles when a context menu event is emitted
	 *
	 * @param action which button was pressed
	 */
	public async handleContextMenu(action: ContextMenuActionEvent<AugmentedChangeLog>): Promise<void> {
		switch (action.action) {
			case "viewChangeLog":
				let changeLog: AugmentedChangeLog;
				if (Array.isArray(action.data)) {
					changeLog = action.data[0];
				} else {
					changeLog = action.data;
				}
				this.dialog.open(TextDialogComponent, {
					data: {message: changeLog.message, title: `Change Log for ${changeLog.user}`}
				});
				break;
		}
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
