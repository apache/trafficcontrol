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

import { Component, OnInit } from "@angular/core";
import { FormControl } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Params } from "@angular/router";
import { BehaviorSubject } from "rxjs";
import type { Region, ResponseRegion } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import type {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink
} from "src/app/shared/generic-table/generic-table.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * RegionsTableComponent is the controller for the "Regions" table.
 */
@Component({
	selector: "tp-regions",
	styleUrls: ["./regions-table.component.scss"],
	templateUrl: "./regions-table.component.html"
})
export class RegionsTableComponent implements OnInit {
	/** List of regions */
	public regions: Promise<Array<ResponseRegion>>;

	constructor(
		private readonly route: ActivatedRoute,
		private readonly navSvc: NavigationService,
		private readonly api: CacheGroupService,
		private readonly dialog: MatDialog,
		public readonly auth: CurrentUserService,
		private readonly log: LoggingService,
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.regions = this.api.getRegions();
		this.navSvc.headerTitle.next("Regions");
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
			},
			e => {
				this.log.error("Failed to get query parameters:", e);
			}
		);
	}

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "name",
			headerName: "Name"
		},
		{
			field: "divisionName",
			headerName: "Division",
			valueFormatter: ({data}: {data: ResponseRegion}): string => `${data.divisionName} (#${data.division})`
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

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseRegion> = {
		href: (row: ResponseRegion): string => `/core/regions/${row.id}`
	};

	/** Definitions for the context menu items (which act on augmented region data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseRegion>> = [
		{
			href: (selectedRow: ResponseRegion): string => `${selectedRow.id}`,
			name: "Edit"
		},
		{
			href: (selectedRow: ResponseRegion): string => `${selectedRow.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		},
		{
			href: (selectedRow: Region): string => `/core/divisions/${selectedRow.division}`,
			name: "View Division"
		},
		{
			href: "/core/phys-locs",
			name: "View Physical Locations",
			queryParams: (selectedRow: ResponseRegion): Params => ({region: selectedRow.name})
		}
	];

	/** A subject that child components can subscribe to for access to the fuzzy search query text */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl = new FormControl<string>("");

	/** Update the URL's 'search' query parameter for the user's search input. */
	public updateURL(): void {
		this.fuzzySubject.next(this.fuzzControl.value ?? "");
	}

	/**
	 * Handles a context menu event.
	 *
	 * @param evt The action selected from the context menu.
	 */
	public async handleContextMenu(evt: ContextMenuActionEvent<ResponseRegion>): Promise<void> {
		const data = evt.data as ResponseRegion;
		switch(evt.action) {
			case "delete":
				const ref = this.dialog.open(DecisionDialogComponent, {
					data: {message: `Are you sure you want to delete region ${data.name} with id ${data.id}`, title: "Confirm Delete"}
				});
				ref.afterClosed().subscribe(result => {
					if(result) {
						this.api.deleteRegion(data.id).then(async () => this.regions = this.api.getRegions());
					}
				});
				break;
		}
	}
}
