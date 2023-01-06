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
import { ActivatedRoute } from "@angular/router";
import type { ColDef } from "ag-grid-community";
import { BehaviorSubject } from "rxjs";
import { LocalizationMethod, localizationMethodToString, ResponseCacheGroup } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import type { ContextMenuActionEvent, ContextMenuItem } from "src/app/shared/generic-table/generic-table.component";
import { TpHeaderService } from "src/app/shared/tp-header/tp-header.service";

/**
 * CacheGroupTableComponent is the controller for the "Cache Groups" table.
 */
@Component({
	selector: "tp-cache-group-table",
	styleUrls: ["./cache-group-table.component.scss"],
	templateUrl: "./cache-group-table.component.html",
})
export class CacheGroupTableComponent implements OnInit {

	/** All of the servers which should appear in the table. */
	public readonly cacheGroups: Promise<Array<ResponseCacheGroup>>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs: ColDef[] = [
		{
			field: "fallbackToClosest",
			filter: "tpBooleanFilter",
			headerName: "Fall-back To Closest",
			hide: false
		},
		{
			field: "id",
			filter: "agNumberColumnFilter",
			headerName: "ID",
			hide: true,
		},
		{
			field: "lastUpdated",
			filter: "agDateColumnFilter",
			headerName: "Last Updated",
			hide: false
		},
		{
			field: "localizationMethods",
			headerName: "Enabled Localization Methods",
			hide: true,
			valueGetter: ({data}: {data: ResponseCacheGroup}): string => {
				let methods;
				if (data.localizationMethods.length > 0) {
					methods = data.localizationMethods;
				} else {
					methods = [LocalizationMethod.CZ, LocalizationMethod.DEEP_CZ, LocalizationMethod.GEO];
				}
				return methods.map(localizationMethodToString).join(", ");
			},
		},
		{
			field: "latitude",
			filter: "agNumberColumnFilter",
			headerName: "Latitude",
			hide: false
		},
		{
			field: "longitude",
			filter: "agNumberColumnFilter",
			headerName: "Longitude",
			hide: false
		},
		{
			field: "name",
			headerName: "Name",
			hide: true,
		},
		{
			field: "parentCachegroupName",
			headerName: "Parent",
			hide: false,
			valueFormatter: (
				{data}: {data: ResponseCacheGroup}
			): string => data.parentCachegroupId === null ? "" : `${data.parentCachegroupName} (#${data.parentCachegroupId})`,
		},
		{
			field: "secondaryParentCachegroupName",
			headerName: "Secondary Parent",
			hide: true,
			valueFormatter: (
				{data}: {data: ResponseCacheGroup}
			): string => data.secondaryParentCachegroupId === null ?
				"" :
				`${data.secondaryParentCachegroupName} (#${data.secondaryParentCachegroupId})`,
		},
		{
			field: "shortName",
			headerName: "Short Name",
			hide: true,
		},
		{
			field: "typeName",
			headerName: "Type",
			hide: false,
			valueFormatter: ({data}: {data: ResponseCacheGroup}): string => `${data.typeName} (#${data.typeId})`
		}
	];

	/**
	 * Definitions for the context menu items (which act on augmented
	 * Cache Group data).
	 */
	public contextMenuItems: Array<ContextMenuItem<ResponseCacheGroup>> = [
		{
			action: "edit",
			name: "Edit"
		},
		{
			action: "delete",
			name: "Delete"
		},
		{
			action: "queue",
			multiRow: true,
			name: "Queue Server Updates"
		},
		{
			action: "dequeue",
			multiRow: true,
			name: "Clear Queued Updates"
		},
		{
			action: "asns",
			name: "Manage ASNs"
		},
		{
			action: "parameters",
			name: "Manage Parameters"
		},
		{
			action: "servers",
			name: "Manage Servers"
		}
	];

	/**
	 * A subject that child components can subscribe to for access to the fuzzy
	 * search query text.
	 */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl: FormControl = new FormControl("");

	constructor(
		private readonly api: CacheGroupService,
		private readonly route: ActivatedRoute,
		private readonly headerSvc: TpHeaderService
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.cacheGroups = this.api.getCacheGroups();
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
				console.error("Failed to get query parameters:", e);
			}
		);
		this.headerSvc.headerTitle.next("Cache Groups");
	}

	/** Update the URL's 'search' query parameter for the user's search input. */
	public updateURL(): void {
		this.fuzzySubject.next(this.fuzzControl.value);
	}

	/**
	 * Handles a context menu event.
	 *
	 * @param a The action selected from the context menu.
	 */
	public handleContextMenu(a: ContextMenuActionEvent<ResponseCacheGroup>): void {
		console.log("action:", a);
	}

}
