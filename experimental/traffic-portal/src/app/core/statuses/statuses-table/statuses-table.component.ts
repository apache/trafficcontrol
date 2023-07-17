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
import { Component, OnInit } from "@angular/core";
import { FormControl } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { BehaviorSubject } from "rxjs";
import { ResponseStatus } from "trafficops-types";

import { ServerService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink
} from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * StatusesTableComponent is the controller for the statuses page - which
 * principally contains a table.
 */
@Component({
	selector: "tp-statuses-table",
	styleUrls: ["./statuses-table.component.scss"],
	templateUrl: "./statuses-table.component.html",
})
export class StatusesTableComponent implements OnInit {

	/** All of the statues which should appear in the table. */
	public statuses: Promise<Array<ResponseStatus>>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "name",
			headerName: "Name",
			hide: false
		},
		{
			field: "description",
			headerName: "Description",
			hide: false
		},
		{
			field: "id",
			headerName: "ID",
			hide: true
		},
		{
			field: "lastUpdated",
			headerName: "Last Updated",
			hide: true
		}];

	/** The current search text. */
	public searchText = "";

	/** Definitions for the context menu items (which act on statuses data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseStatus>> = [
		{
			href: (status: ResponseStatus): string => `${status.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			href: (status: ResponseStatus): string => `${status.id}`,
			name: "Edit"
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		}
	];

	/** A subject that child components can subscribe to for access to the fuzzy search query text */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl = new FormControl<string>("", {nonNullable: true});

	constructor(
		private readonly dialog: MatDialog,
		private readonly route: ActivatedRoute,
		private readonly api: ServerService,
		private readonly navSvc: NavigationService, public readonly auth: CurrentUserService
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.statuses = this.api.getStatuses();
		this.navSvc.headerTitle.next("Statuses");
	}

	/** Initializes table data, loading it from Traffic Ops. */
	public ngOnInit(): void {
		this.route.queryParamMap.subscribe(
			m => {
				const search = m.get("search");
				if (search) {
					this.fuzzControl.setValue(decodeURIComponent(search));
					this.updateURL();
				}
			}
		);
	}

	/**
	 * Updates the "search" query parameter in the URL every time the search
	 * text input changes.
	 */
	public updateURL(): void {
		this.fuzzySubject.next(this.searchText);
	}

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseStatus> = {
		href: (row: ResponseStatus): string => `/core/statuses/${row.id}`
	};

	/**
	 * Handles a context menu event.
	 *
	 * @param evt The action selected from the context menu.
	 */
	public async handleContextMenu(evt: ContextMenuActionEvent<ResponseStatus>): Promise<void> {
		const data = evt.data as ResponseStatus;
		switch(evt.action) {
			case "delete":
				const ref = this.dialog.open(DecisionDialogComponent, {
					data: {message: `Are you sure you want to delete Status ${data.name} with ID ${data.id} ?`, title: "Confirm Delete"}
				});
				ref.afterClosed().subscribe(result => {
					if(result) {
						this.api.deleteStatus(data.id).then(async () => this.statuses = this.api.getStatuses());
					}
				});
				break;
		}
	}
}
