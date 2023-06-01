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
import { ActivatedRoute, Params } from "@angular/router";
import { BehaviorSubject } from "rxjs";
import { ResponseParameter } from "trafficops-types";

import { ProfileService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import type { ContextMenuActionEvent, ContextMenuItem, DoubleClickLink } from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * ParametersTableComponent is the controller for the "Parameters" table.
 */
@Component({
	selector: "tp-parameters",
	styleUrls: ["./parameters-table.component.scss"],
	templateUrl: "./parameters-table.component.html"
})
export class ParametersTableComponent implements OnInit {
	/** List of parameters */
	public parameters: Promise<Array<ResponseParameter>>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "id",
			filter: "agNumberColumnFilter",
			headerName: "ID",
			hide: true
		},
		{
			field: "configFile",
			headerName: "Config File"
		},
		{
			field: "name",
			headerName: "Name"
		},
		{
			field: "profiles",
			headerName: "Profiles",
			valueFormatter: ({data}: {data: ResponseParameter}): string => data.profiles === null? "":data.profiles.join(", ")
		},
		{
			field: "secure",
			filter: "tpBooleanFilter",
			headerName: "Secure",
			hide: true
		},
		{
			field: "value",
			headerName: "Value"
		},
		{
			field: "lastUpdated",
			filter: "agDateColumnFilter",
			headerName: "Last Updated",
			hide: true
		}
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseParameter> = {
		href: (row: ResponseParameter): string => `/core/parameters/${row.id}`
	};

	/** Definitions for the context menu items (which act on augmented parameter data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseParameter>> = [
		{
			href: (responseParameter: ResponseParameter): string => `${responseParameter.id}`,
			name: "Edit"
		},
		{
			href: (responseParameter: ResponseParameter): string => `${responseParameter.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		},
		{
			href: "/core/profiles",
			name: "View Profiles",
			queryParams: (selectedRow: ResponseParameter): Params => ({hasParameter: selectedRow.id}),
		}
	];

	/** A subject that child components can subscribe to for access to the fuzzy search query text */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl = new FormControl<string>("", {nonNullable: true});

	constructor(private readonly route: ActivatedRoute, private readonly navSvc: NavigationService,
		private readonly api: ProfileService, private readonly dialog: MatDialog, public readonly auth: CurrentUserService) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.parameters = this.api.getParameters();
		this.navSvc.headerTitle.next("Parameters");
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
	public async handleContextMenu(evt: ContextMenuActionEvent<ResponseParameter>): Promise<void> {
		const data = evt.data as ResponseParameter;
		switch(evt.action) {
			case "delete":
				const ref = this.dialog.open(DecisionDialogComponent, {
					data: {message: `Are you sure you want to delete Parameter ${data.name} with ID ${data.id}?`, title: "Confirm Delete"}
				});
				ref.afterClosed().subscribe(result => {
					if(result) {
						this.api.deleteParameter(data.id).then(async () => this.parameters = this.api.getParameters());
					}
				});
				break;
		}
	}
}
