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
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { ResponseDivision } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { AbstractTableComponent } from "src/app/shared/generic-table/abstract-table.component";
import type { ContextMenuActionEvent, ContextMenuItem } from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * DivisionsTableComponent is the controller for the "Divisions" table.
 */
@Component({
	selector: "tp-divisions",
	styleUrls: ["../../../../shared/generic-table/abstract-table.component.scss"],
	templateUrl: "../../../../shared/generic-table/abstract-table.component.html",
})
export class DivisionsTableComponent extends AbstractTableComponent<ResponseDivision> implements OnInit {
	public readonly context = "divisions";
	public readonly tableName = "Divisions";

	public override readonly fabTitle = "Create a new Division";
	public override readonly fabLink = "new";
	public override readonly fabType = "link";

	/** List of divisions */
	public data: Promise<Array<ResponseDivision>>;

	constructor(
		route: ActivatedRoute,
		navSvc: NavigationService,
		private readonly api: CacheGroupService,
		private readonly dialog: MatDialog,
		public readonly auth: CurrentUserService
	) {
		super(route, navSvc);
		this.data = this.api.getDivisions();
	}

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
			href: (div: ResponseDivision): string => `core/divisions/${div.id}`,
			name: "Edit"
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		},
		{
			href: (div: ResponseDivision): string => `core/regions?search=${div.name}`,
			name: "View Regions"
		}
	];

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
						this.api.deleteDivision(data.id).then(async () => this.data = this.api.getDivisions());
					}
				});
				break;
		}
	}

	/**
	 * Checks if the user has permission to use the Division table FAB.
	 *
	 * @returns `true` if the user has permission to create Divisions, `false`
	 * otherwise.
	 */
	public override fabPermission(): boolean {
		return this.auth.hasPermission("DIVISION:CREATE");
	}
}
