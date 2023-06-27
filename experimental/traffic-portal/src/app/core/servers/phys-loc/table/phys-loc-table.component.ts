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
import { Component, type OnInit } from "@angular/core";
import { FormControl } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, type Params } from "@angular/router";
import { BehaviorSubject } from "rxjs";
import type { ResponsePhysicalLocation } from "trafficops-types";

import { PhysicalLocationService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import type {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink
} from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * PHysLocTableComponent is the controller for the phys loc table.
 */
@Component({
	selector: "tp-phys-loc-table",
	styleUrls: ["./phys-loc-table.component.scss"],
	templateUrl: "./phys-loc-table.component.html"
})
export class PhysLocTableComponent implements OnInit {
	/** All the physical locations which should appear in the table. */
	public physLocations: Promise<Array<ResponsePhysicalLocation>>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [{
		field: "name",
		headerName: "Name"
	}, {
		field: "id",
		headerName: "ID",
		hide: true
	}, {
		field: "shortName",
		headerName: "Short Name"
	}, {
		field: "address",
		headerName: "Address"
	}, {
		field: "city",
		headerName: "City"
	}, {
		field: "state",
		headerName: "State"
	}, {
		field: "region",
		headerName: "Region",
		valueFormatter: ({data}: {data: ResponsePhysicalLocation}): string => `${data.region} (#${data.regionId})`
	}, {
		field: "lastUpdated",
		headerName: "Last Updated"
	}];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponsePhysicalLocation> = {
		href: (row: ResponsePhysicalLocation): string => `/core/phys-locs/${row.id}`
	};

	/** Definitions for the context menu items (which act on augmented cache-group data). */
	public contextMenuItems: Array<ContextMenuItem<ResponsePhysicalLocation>> = [
		{
			href: (physLoc: ResponsePhysicalLocation): string => `${physLoc.id}`,
			name: "Edit"
		},
		{
			href: (physLoc: ResponsePhysicalLocation): string => `${physLoc.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		},
		{
			href: (physLoc: ResponsePhysicalLocation): string => `/core/regions/${physLoc.regionId}`,
			name: "View Region"
		},
		{
			href: "/core/servers",
			name: "View Servers",
			queryParams: (physLoc: ResponsePhysicalLocation): Params => ({physLocation: physLoc.name})
		}
	];
	/** A subject that child components can subscribe to for access to the fuzzy search query text */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl: FormControl = new FormControl<string>("");

	constructor(
		private readonly api: PhysicalLocationService,
		private readonly route: ActivatedRoute,
		private readonly navSvc: NavigationService,
		public readonly auth: CurrentUserService,
		private readonly dialog: MatDialog,
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.physLocations = this.api.getPhysicalLocations();
		this.navSvc.headerTitle.next("Physical Locations");
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
	public async handleContextMenu(evt: ContextMenuActionEvent<ResponsePhysicalLocation>): Promise<void> {
		const data = evt.data as ResponsePhysicalLocation;
		switch(evt.action) {
			case "delete":
				const ref = this.dialog.open(DecisionDialogComponent, {
					data: {
						message: `Are you sure you want to delete Physical Location ${data.name} with id ${data.id}?`,
						title: "Confirm Delete"
					}
				});
				ref.afterClosed().subscribe(result => {
					if(result) {
						this.api.deletePhysicalLocation(data.id).then(
							async () => this.physLocations = this.api.getPhysicalLocations()
						);
					}
				});
				break;
		}
	}

}
