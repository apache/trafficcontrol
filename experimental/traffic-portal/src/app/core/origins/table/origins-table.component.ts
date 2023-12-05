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
import { ActivatedRoute } from "@angular/router";
import { BehaviorSubject } from "rxjs";
import { TypeFromResponse } from "trafficops-types";
import { ResponseOrigin } from "trafficops-types/dist/origin";

import { OriginService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import type {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink
} from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * OriginsTableComponent is the controller for the "Origins" table.
 */
@Component({
	selector: "tp-origins",
	styleUrls: ["./origins-table.component.scss"],
	templateUrl: "./origins-table.component.html"
})
export class OriginsTableComponent implements OnInit {
	/** List of origins */
	public origins: Promise<Array<ResponseOrigin>>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "name",
			headerName: "Name"
		},
		{
			field: "tenant",
			headerName: "Tenant"
		},
		{
			field: "primary",
			headerName: "Primary"
		},
		{
			field: "delivery-service",
			headerName: "Delivery Service"
		},
		{
			field: "fqdn",
			headerName: "FQDN"
		},
		{
			field: "lastUpdated",
			headerName: "Last Updated"
		}
	];

	/** Definitions for the context menu items (which act on augmented type data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseOrigin>> = [
		{
			href: (type: ResponseOrigin): string => `${type.id}`,
			name: "Edit"
		},
		{
			href: (type: ResponseOrigin): string => `${type.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		}
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseOrigin> = {
		href: (row: ResponseOrigin): string => `/core/origins/${row.id}`
	};

	/** A subject that child components can subscribe to for access to the fuzzy search query text */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl = new FormControl<string>("", {nonNullable: true});

	constructor(private readonly route: ActivatedRoute, private readonly navSvc: NavigationService,
		private readonly api: OriginService, private readonly dialog: MatDialog, public readonly auth: CurrentUserService) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.origins = this.api.getOrigins();
		this.navSvc.headerTitle.next("Origins");
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
	public async handleContextMenu(evt: ContextMenuActionEvent<TypeFromResponse>): Promise<void> {
		const data = evt.data as TypeFromResponse;
		switch(evt.action) {
			case "delete":
				const ref = this.dialog.open(DecisionDialogComponent, {
					data: {message: `Are you sure you want to delete origin ${data.name} with id ${data.id} ?`, title: "Confirm Delete"}
				});
				ref.afterClosed().subscribe(result => {
					if(result) {
						this.api.deleteOrigin(data.id).then(async () => this.origins = this.api.getOrigins());
					}
				});
				break;
		}
	}
}
