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
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { Region, ResponseRegion } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { AbstractTableComponent } from "src/app/shared/generic-table/abstract-table.component";
import { ContextMenuActionEvent, ContextMenuItem } from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * RegionsTableComponent is the controller for the "Regions" table.
 */
@Component({
	selector: "tp-regions",
	styleUrls: ["../../../../shared/generic-table/abstract-table.component.scss"],
	templateUrl: "../../../../shared/generic-table/abstract-table.component.html",
})
export class RegionsTableComponent extends AbstractTableComponent<ResponseRegion> implements OnInit {
	/** List of regions */
	public data: Promise<Array<ResponseRegion>>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "name",
			headerName: "Name"
		},
		{
			field: "divisionName",
			headerName: "Division"
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

	public readonly context = "regions";
	public readonly tableName = "Regions";

	public override readonly fabTitle = "Add new Region";
	public override readonly fabLink = "new";
	public override readonly fabType = "link";

	constructor(
		route: ActivatedRoute,
		headerSvc: NavigationService,
		private readonly api: CacheGroupService,
		private readonly dialog: MatDialog,
		public readonly auth: CurrentUserService
	) {
		super(route, headerSvc);
		this.data = this.api.getRegions();
	}

	/** Definitions for the context menu items (which act on augmented region data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseRegion>> = [
		{
			href: (selectedRow: ResponseRegion): string => `/core/regions/${selectedRow.id}`,
			name: "Edit"
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
			href: (selectedRow: Region): string => `/core/phys-locs?search=${selectedRow.name}`,
			name: "View Physical Locations"
		}
	];

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
						this.api.deleteRegion(data.id).then(async () => this.data = this.api.getRegions());
					}
				});
				break;
		}
	}

	/**
	 * Checks if the user has permission to use the Regions table FAB.
	 *
	 * @returns `true` if the user has permission to create Regions, `false`
	 * otherwise.
	 */
	public override fabPermission(): boolean {
		return this.auth.hasPermission("REGION:CREATE");
	}
}
