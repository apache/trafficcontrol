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
import { ActivatedRoute } from "@angular/router";
import { BehaviorSubject } from "rxjs";
import type { ResponseServerCapability } from "trafficops-types";

import { ServerService } from "src/app/api";
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
 * Controller for the table that displays Server Capabilities.
 */
@Component({
	selector: "tp-capabilities",
	styleUrls: ["./capabilities.component.scss"],
	templateUrl: "./capabilities.component.html",
})
export class CapabilitiesComponent implements OnInit {
	/** All the physical locations which should appear in the table. */
	public capabilities: Promise<Array<ResponseServerCapability>>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "name",
			headerName: "Name"
		},
		{
			field: "lastUpdated",
			headerName: "Last Updated",
			hide: true
		},
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseServerCapability> = {
		href: (row: ResponseServerCapability): string => `/core/capabilities/${row.name}`
	};

	/** Definitions for the context menu items (which act on augmented cache-group data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseServerCapability>> = [
		{
			href: (c: ResponseServerCapability): string => c.name,
			name: "View Details",
		},
		{
			href: (c: ResponseServerCapability): string => c.name,
			name: "Open in New Tab",
			newTab: true,
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		},
		{
			action: "servers",
			// TODO: implement
			disabled: (): true => true,
			multiRow: true,
			name: "View Servers",
		},
		{
			action: "addServers",
			// TODO: implement
			disabled: (): true => true,
			multiRow: true,
			name: "Add to Server(s)",
		},
		{
			action: "dses",
			// TODO: implement
			disabled: (): true => true,
			multiRow: true,
			name: "View Delivery Services",
		},
		{
			action: "addDSes",
			// TODO: implement
			disabled: (): true => true,
			multiRow: true,
			name: "Add to Delivery Service(s)",
		},
	];

	/** A subject that child components can subscribe to for access to the fuzzy search query text */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl = new FormControl<string>("", {nonNullable: true});

	/**
	 * Constructs the component with its required injections.
	 *
	 * @param api The Servers API which is used to provide row data.
	 * @param route A reference to the route of this view which is used to set the fuzzy search box text from the 'search' query parameter.
	 * @param router Angular router
	 * @param navSvc Manages the header
	 * @param dialog Dialog manager
	 */
	constructor(
		private readonly api: ServerService,
		private readonly route: ActivatedRoute,
		private readonly navSvc: NavigationService,
		private readonly dialog: MatDialog,
		public readonly auth: CurrentUserService,
		private readonly log: LoggingService
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.capabilities = this.api.getCapabilities();
		this.navSvc.headerTitle.next("Capabilities");
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
	public async handleContextMenu(evt: ContextMenuActionEvent<ResponseServerCapability>): Promise<void> {
		const data = evt.data;
		switch (evt.action) {
			case "delete":
				if (Array.isArray(data)) {
					throw new Error("cannot delete multiple Capabilities");
				}
				const ref = this.dialog.open(DecisionDialogComponent, {
					data: {
						message: `Are you sure you want to delete the '${data.name}' Capability?`,
						title: "Confirm Delete"
					}
				});
				const result = await ref.afterClosed().toPromise();
				if (result) {
					this.api.deleteCapability(data).then(async () => this.capabilities = this.api.getCapabilities());
				}
				break;
			default:
				this.log.warn("unrecognized context menu action:", evt.action);
		}
	}
}
