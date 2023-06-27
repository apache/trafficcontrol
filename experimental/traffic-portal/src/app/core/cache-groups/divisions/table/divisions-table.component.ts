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

import { Component, type OnInit } from "@angular/core";
import { FormControl } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, type Params } from "@angular/router";
import { BehaviorSubject } from "rxjs";
import { ResponseDivision } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import type {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink
} from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * DivisionsTableComponent is the controller for the "Divisions" table.
 */
@Component({
	selector: "tp-divisions",
	styleUrls: ["./divisions-table.component.scss"],
	templateUrl: "./divisions-table.component.html"
})
export class DivisionsTableComponent implements OnInit {
	/** List of divisions */
	public divisions: Promise<Array<ResponseDivision>>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "name",
			headerName: "Name"
		},
		{
			field: "id",
			headerName:" ID",
			hide: true
		},
		{
			field: "lastUpdated",
			headerName: "Last Updated"
		}
	];

	/** Definitions for the context menu items (which act on augmented division data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseDivision>> = [
		{
			href: (div: ResponseDivision): string => `${div.id}`,
			name: "Edit"
		},
		{
			href: (div: ResponseDivision): string => `${div.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		},
		{
			href: "/core/regions",
			name: "View Regions",
			queryParams: (div: ResponseDivision): Params => ({divisionName: div.name})
		}
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseDivision> = {
		href: (row: ResponseDivision): string => `/core/divisions/${row.id}`
	};

	/** A subject that child components can subscribe to for access to the fuzzy search query text */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl = new FormControl<string>("", {nonNullable: true});

	constructor(private readonly route: ActivatedRoute, private readonly navSvc: NavigationService,
		private readonly api: CacheGroupService, private readonly dialog: MatDialog, public readonly auth: CurrentUserService) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.divisions = this.api.getDivisions();
		this.navSvc.headerTitle.next("Divisions");
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

	/** Update the URL's 'search' query parameter for the user's search input. */
	public updateURL(): void {
		this.fuzzySubject.next(this.fuzzControl.value);
	}

	/**
	 * Handles a context menu event.
	 *
	 * @param evt The action selected from the context menu.
	 */
	public async handleContextMenu(evt: ContextMenuActionEvent<ResponseDivision>): Promise<void> {
		const data = evt.data as ResponseDivision;
		switch(evt.action) {
			case "delete":
				const ref = this.dialog.open(DecisionDialogComponent, {
					data: {message: `Are you sure you want to delete division ${data.name} with id ${data.id}`, title: "Confirm Delete"}
				});
				ref.afterClosed().subscribe(result => {
					if(result) {
						this.api.deleteDivision(data.id).then(async () => this.divisions = this.api.getDivisions());
					}
				});
				break;
		}
	}
}
