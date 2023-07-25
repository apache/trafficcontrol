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

import { Component, type OnDestroy, type OnInit } from "@angular/core";
import type { Params } from "@angular/router";
import type { ValueFormatterParams } from "ag-grid-community";
import { BehaviorSubject, type Subscription } from "rxjs";
import { ResponseTenant } from "trafficops-types";

import { UserService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import type {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink
} from "src/app/shared/generic-table/generic-table.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * TenantsComponent is the controller for the table that lists Tenants.
 */
@Component({
	selector: "tp-tenants",
	styleUrls: ["./tenants.component.scss"],
	templateUrl: "./tenants.component.html"
})
export class TenantsComponent implements OnInit, OnDestroy {

	private tenantMap: Record<number, ResponseTenant> = {};

	public searchText = "";
	public searchSubject = new BehaviorSubject("");

	public tenants: Array<ResponseTenant> = [{
		active: true,
		id: 1,
		lastUpdated: new Date(),
		name: "root",
		parentId: null
	}];

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "active",
			filter: "tpBooleanFilter",
			headerName: "Active",
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
			hide: true,
		},
		{
			field: "name",
			headerName: "Name",
			hide: false,
		},
		{
			field: "parentId",
			headerName: "Parent",
			hide: false,
			valueFormatter: (params: ValueFormatterParams): string => this.getParentString(params.data)
		}
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseTenant> = {
		href: (row: ResponseTenant): string => `/core/tenants/${row.id}`
	};

	public contextMenuItems: ContextMenuItem<Readonly<ResponseTenant>>[] = [];

	public loading = true;
	private readonly subscription: Subscription;

	constructor(
		private readonly userService: UserService,
		public readonly auth: CurrentUserService,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) {
		this.navSvc.headerTitle.next("Tenant");
		this.subscription = this.auth.userChanged.subscribe(
			() => {
				this.loadContextMenuItems();
			}
		);
	}

	/**
	 * Loads the context menu items for the grid.
	 *
	 * @private
	 */
	private loadContextMenuItems(): void {
		this.contextMenuItems = [];
		if (this.auth.hasPermission("TENANT:UPDATE")) {
			this.contextMenuItems.push({
				href: (t: ResponseTenant): string => `${t.id}`,
				name: "View Details"
			});
			this.contextMenuItems.push({
				href: (t: ResponseTenant): string => `${t.id}`,
				name: "Open in New Tab",
				newTab: true
			});
			this.contextMenuItems.push({
				action: "disable",
				disabled: (ts): boolean => ts.some(t => t.name === "root" || t.id === this.auth.currentUser?.tenantId),
				multiRow: true,
				name: "Disable"
			});
		}
		this.contextMenuItems.push({
			disabled: (t: ResponseTenant | ResponseTenant[]): boolean =>
				Array.isArray(t) || t.id === this.auth.currentUser?.tenantId || t.parentId === null,
			href: (t: ResponseTenant): string => `${t.parentId}`,
			name: "View Parent Details"
		});
		if (this.auth.hasPermission("USER:READ")) {
			this.contextMenuItems.push({
				href: "/core/users",
				name: "View Users",
				queryParams: (t: ResponseTenant): Params => ({tenant: t.name})
			});
		}
	}

	/**
	 * Angular lifecycle hook; fetches API data.
	 */
	public async ngOnInit(): Promise<void> {
		this.tenants = await this.userService.getTenants();
		this.tenantMap = Object.fromEntries((this.tenants).map(t => [t.id, t]));
		this.loadContextMenuItems();
		this.loading = false;
	}

	/**
	 * Gets a string representation for the Parent of the given Tenant.
	 *
	 * @param t The Tenant for which the Parent will be rendered.
	 * @returns An empty string for the root Tenant, otherwise the parent
	 * Tenant's name and ID as a string.
	 */
	public getParentString(t: ResponseTenant): string {
		if (t.parentId === null) {
			return "";
		}
		return `${this.tenantMap[t.parentId].name} (#${t.parentId})`;
	}

	/**
	 * Angular lifecycle hook; cleans up persistent resources.
	 */
	public ngOnDestroy(): void {
		this.subscription.unsubscribe();
	}

	/**
	 * Handles a context menu event.
	 *
	 * @param a The action selected from the context menu.
	 */
	public handleContextMenu(a: ContextMenuActionEvent<Readonly<ResponseTenant>>): void {
		this.log.debug("action:", a);
	}

	/**
	 * Updates the "search" query parameter in the URL every time the search
	 * text input changes.
	 */
	public updateURL(): void {
		this.searchSubject.next(this.searchText);
	}
}
