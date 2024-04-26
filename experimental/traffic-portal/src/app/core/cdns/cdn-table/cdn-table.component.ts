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
import { ColDef } from "ag-grid-community";
import { BehaviorSubject } from "rxjs";
import { AlertLevel, ResponseCDN } from "trafficops-types";

import { CDNService } from "src/app/api";
import { AlertService } from "src/app/shared/alert/alert.service";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import {
	DecisionDialogComponent,
	DecisionDialogData
} from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import type {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink
} from "src/app/shared/generic-table/generic-table.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * CDNTableComponent is the controller for the "CDNs" table.
 */
@Component({
	selector: "tp-cdn-table",
	styleUrls: ["./cdn-table.component.scss"],
	templateUrl: "./cdn-table.component.html",
})
export class CDNTableComponent implements OnInit {
	public cdns: Promise<ResponseCDN[]>;

	/* Definitions of the table's columns according to the ag-grid API */
	public columnDefs: ColDef[] = [
		{
			field: "dnssecEnabled",
			filter: "tpBooleanFilter",
			headerName: "DNSSEC Enabled",
			hide: false
		},
		{
			field: "domainName",
			filter: "agTextColumnFilter",
			headerName: "Domain",
			hide: false,
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
			hide: true
		},
		{
			field: "name",
			filter: "agTextColumnFilter",
			headerName: "Name",
			hide: false,
		},
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseCDN> = {
		href: (row: ResponseCDN): string => `/core/cdns/${row.id}`
	};

	/**
	 * Definitions for the context menu items (which act on augmented
	 * CDN data).
	 */
	public contextMenuItems: Array<ContextMenuItem<ResponseCDN>> = [
		{
			href: (selectedRow): string => `/core/cdns/${selectedRow.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			href: (selectedRow): string => `/core/cdns/${selectedRow.id}`,
			name: "Edit"
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete",
		},
		{
			action: "snapshot-diff",
			disabled: (): true => true,
			multiRow: false,
			name: "Diff Snapshot",
		},
		{
			action: "queue",
			multiRow: false,
			name: "Queue Server Updates"
		},
		{
			action: "dequeue",
			multiRow: false,
			name: "Clear Queued Updates"
		},
		{
			disabled: (): true => true,
			href: (selectedRow): string => `/core/cdns/${selectedRow.id}/dnssec-keys`,
			name: "Manage DNSSEC Keys"
		},
		{
			disabled: (): true => true,
			href: (selectedRow): string => `/core/cdns/${selectedRow.id}/federations`,
			name: "Manage Federations"
		},
		{
			disabled: (): true => true,
			href: (selectedRow): string => `/core/cdns/${selectedRow.id}/delivery-services`,
			name: "Manage Delivery Services"
		},
		{
			href: "/core/profiles",
			name: "View Profiles",
			queryParams: (selectedRow): Params => ({cdnName: selectedRow.name}),
		},
		{
			href: "/core/servers",
			name: "View Servers",
			queryParams: (selectedRow): Params => ({cdnName: selectedRow.name}),
		},
		{
			disabled: (): true => true,
			href: (selectedRow): string => `/core/cdns/${selectedRow.id}/notifications`,
			name: "Manage Notifications"
		},
	];

	/**
	 * A subject that child components can subscribe to for access to the fuzzy
	 * search query text.
	 */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl: FormControl = new FormControl("");

	constructor(
		private readonly alerts: AlertService,
		private readonly api: CDNService,
		public readonly auth: CurrentUserService,
		private readonly navSvc: NavigationService,
		private readonly dialog: MatDialog,
		private readonly route: ActivatedRoute,
		private readonly log: LoggingService,
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.cdns = this.api.getCDNs();
		this.navSvc.headerTitle.next("CDNs");
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
	 * Queues or clears updates on a group of CDNs.
	 *
	 * @param cdn The CDN on which to queue updates.
	 * @param queue Whether updates should be queued (`true`) or cleared
	 * (`false`).
	 */
	private async queueUpdates(cdn: ResponseCDN | number, queue: boolean = true): Promise<void> {
		if (typeof cdn === "number") {
			cdn = await this.api.getCDNs(cdn);
		}
		const readableAction = queue ? "Queue" : "Clear";
		const title = `${readableAction} Updates on ${cdn.name}?`;
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(DecisionDialogComponent, {
			data: {
				message: `Are you sure you want to ${readableAction.toLowerCase()} server updates for all of the ${cdn.name} servers?`,
				title,
			}
		});
		if (!await ref.afterClosed().toPromise()) {
			return;
		}
		if (queue) {
			await this.api.queueServerUpdates(cdn);
		} else {
			await this.api.dequeueServerUpdates(cdn);
		}
		this.alerts.newAlert(
			AlertLevel.SUCCESS,
			`${readableAction.replace(/(?<!e)$/, "e")}d CDN server updates`,
		);
	}

	/**
	 * Asks the user for confirmation before deleting a CDN.
	 *
	 * @param cdn The CDN (potentially) being deleted.
	 */
	private async delete(cdn: ResponseCDN | number): Promise<void> {
		if (typeof cdn === "number") {
			cdn = await this.api.getCDNs(cdn);
		}
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(DecisionDialogComponent, {
			data: {
				message: `Are you sure you want to delete the ${cdn.name} CDN?`,
				title: `Delete ${cdn.name}`
			}
		});
		if (await ref.afterClosed().toPromise()) {
			await this.api.deleteCDN(cdn);
			this.cdns = this.api.getCDNs();
		}
	}

	/**
	 * Handles a context menu event.
	 *
	 * @param a The action selected from the context menu.
	 */
	public handleContextMenu(a: ContextMenuActionEvent<ResponseCDN>): void {
		switch (a.action) {
			case "queue":
				if (Array.isArray(a.data)) {
					this.log.error("cannot queue multiple cdns at once:", a.data);
					return;
				}
				this.queueUpdates(a.data);
				break;
			case "dequeue":
				if (Array.isArray(a.data)) {
					this.log.error("cannot dequeue multiple cdns at once:", a.data);
					return;
				}
				this.queueUpdates(a.data, false);
				break;
			case "delete":
				if (Array.isArray(a.data)) {
					this.log.error("cannot delete multiple cdns at once:", a.data);
					return;
				}
				this.delete(a.data);
				break;
			default:
				this.log.error("unrecognized context menu action:", a.action);
		}
	}
}
