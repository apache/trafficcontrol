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
import type { ColDef } from "ag-grid-community";
import { BehaviorSubject } from "rxjs";
import {
	AlertLevel,
	LocalizationMethod,
	localizationMethodToString,
	type ResponseCacheGroup,
	type ResponseCDN
} from "trafficops-types";

import { CacheGroupService, CDNService } from "src/app/api";
import { AlertService } from "src/app/shared/alert/alert.service";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import {
	CollectionChoiceDialogComponent,
	type CollectionChoiceDialogData
} from "src/app/shared/dialogs/collection-choice-dialog/collection-choice-dialog.component";
import { DecisionDialogComponent, type DecisionDialogData } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import type {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink
} from "src/app/shared/generic-table/generic-table.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

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
	public cacheGroups: Promise<Array<ResponseCacheGroup>>;

	/** All of the CDNs (on which a user might (de/)queue updates). */
	public readonly cdns: Promise<Array<ResponseCDN>>;

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
			href: (selectedRow): string => String(selectedRow.id),
			name: "Open in New Tab",
			newTab: true
		},
		{
			href: (selectedRow): string => String(selectedRow.id),
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
			href: "/core/asns",
			name: "View ASNs",
			queryParams: (selectedRow):  Params => ({cachegroup: selectedRow.name})
		},
		{
			href: "/core/servers",
			name: "View Servers",
			queryParams: (selectedRow): Params => ({cachegroup: selectedRow.name})
		}
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseCacheGroup> = {
		href: (row: ResponseCacheGroup): string => `/core/cache-groups/${row.id}`
	};

	/**
	 * A subject that child components can subscribe to for access to the fuzzy
	 * search query text.
	 */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl: FormControl = new FormControl("");

	constructor(
		private readonly api: CacheGroupService,
		private readonly cdnAPI: CDNService,
		private readonly route: ActivatedRoute,
		private readonly dialog: MatDialog,
		private readonly alerts: AlertService,
		public readonly auth: CurrentUserService,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.cacheGroups = this.api.getCacheGroups();
		this.navSvc.headerTitle.next("Cache Groups");
		this.cdns = this.cdnAPI.getCDNs();
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
	 * Queues or clears updates on a group of Cache Groups.
	 *
	 * @param cgs The Cache Groups on which to operate.
	 * @param queue Whether updates should be queued (`true`) or cleared
	 * (`false`).
	 */
	private async queueUpdates(cgs: ResponseCacheGroup[], queue: boolean = true): Promise<void> {
		const title = `${queue ? "Queue" : "Clear"} Updates on ${cgs.length === 1 ? cgs[0].name : `${cgs.length} Cache Groups`}`;
		const data = {
			collection: (await this.cdns).map(c => ({label: c.name, value: c.id})),
			hint: "Note that 'ALL' does NOT mean 'all CDNs'!",
			label: "CDN",
			message: `Select a CDN to which to limit the ${queue ? "Queuing" : "Clearing"} of Updates.`,
			title,
		};
		const ref = this.dialog.open<CollectionChoiceDialogComponent, CollectionChoiceDialogData<number>, number | false>(
			CollectionChoiceDialogComponent,
			{data}
		);
		const result = await ref.afterClosed().toPromise();
		if (typeof(result) === "number") {
			const responses = await Promise.all(cgs.map(async cg => this.api.queueCacheGroupUpdates(cg, result)));
			const serverNum = responses.map(r => r.serverNames.length).reduce((n, l) => n+l, 0);
			// This endpoint returns no alerts at the time of this writing, so
			// we gotta do it by hand.
			this.alerts.newAlert(
				AlertLevel.SUCCESS,
				`${queue ? "Queued" : "Cleared"} Updates on ${serverNum} server${serverNum === 1 ? "" : "s"}`
			);
		}
	}

	/**
	 * Asks the user for confirmation before deleting a Cache Group.
	 *
	 * @param cg The Cache Group (potentially) being deleted.
	 */
	private async delete(cg: ResponseCacheGroup): Promise<void> {
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(DecisionDialogComponent, {
			data: {
				message: `Are you sure you want to delete the ${cg.name} Cache Group?`,
				title: `Delete ${cg.name}`
			}
		});
		if (await ref.afterClosed().toPromise()) {
			await this.api.deleteCacheGroup(cg);
			this.cacheGroups = this.api.getCacheGroups();
		}
	}

	/**
	 * Handles a context menu event.
	 *
	 * @param a The action selected from the context menu.
	 */
	public handleContextMenu(a: ContextMenuActionEvent<ResponseCacheGroup>): void {
		switch(a.action) {
			case "queue":
				this.queueUpdates(Array.isArray(a.data) ? a.data : [a.data]);
				break;
			case "dequeue":
				this.queueUpdates(Array.isArray(a.data) ? a.data : [a.data], false);
				break;
			case "delete":
				if (Array.isArray(a.data)) {
					this.log.error("cannot delete multiple cache groups at once:", a.data);
					return;
				}
				this.delete(a.data);
				break;
			default:
				this.log.error("unrecognized context menu action:", a.action);
		}
	}
}
