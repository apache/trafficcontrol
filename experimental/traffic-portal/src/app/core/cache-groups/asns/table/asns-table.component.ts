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
import type { ResponseASN } from "trafficops-types";

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
 * AsnsTableComponent is the controller for the "Asns" table.
 */
@Component({
	selector: "tp-asns",
	styleUrls: ["./asns-table.component.scss"],
	templateUrl: "./asns-table.component.html"
})
export class ASNsTableComponent implements OnInit {
	/** List of asns */
	public asns: Promise<Array<ResponseASN>>;

	constructor(
		private readonly route: ActivatedRoute,
		private readonly navSvc: NavigationService,
		private readonly api: CacheGroupService,
		private readonly dialog: MatDialog,
		public readonly auth: CurrentUserService,
		private readonly log: LoggingService,
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.asns = this.api.getASNs();
		this.navSvc.headerTitle.next("ASNs");
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
			field: "asn",
			headerName: "ASN"
		},
		{
			field: "cachegroup",
			headerName: "Cache Group",
		},
		{
			field: "lastUpdated",
			headerName: "Last Updated"
		}
	];

	/** Definitions for the context menu items (which act on augmented asn data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseASN>> = [
		{
			href: (selectedRow: ResponseASN): string => `${selectedRow.id}`,
			name: "Edit"
		},
		{
			href: (selectedRow: ResponseASN): string => `${selectedRow.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		},
		{
			href: "/core/cache-groups",
			name: "View Cache Group",
			queryParams: (selectedRow: ResponseASN): Params => ({name: selectedRow.cachegroup}),
		}
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseASN> = {
		href: (row: ResponseASN): string => `/core/asns/${row.id}`
	};

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
	public async handleContextMenu(evt: ContextMenuActionEvent<ResponseASN>): Promise<void> {
		if (Array.isArray(evt.data)) {
			this.log.error("cannot delete multiple ASNs at once:", evt.data);
			return;
		}
		const data = evt.data;
		switch(evt.action) {
			case "delete":
				const ref = this.dialog.open(DecisionDialogComponent, {
					data: {message: `Are you sure you want to delete ASN ${data.asn}?`, title: "Confirm Delete"}
				});
				ref.afterClosed().subscribe(result => {
					if(result) {
						this.api.deleteASN(data.id).then(async () => this.asns = this.api.getASNs());
					}
				});
				break;
		}
	}
}
