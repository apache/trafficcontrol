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
import type { ValueGetterParams } from "ag-grid-community";
import { BehaviorSubject, type Subscription } from "rxjs";

import { UserService } from "src/app/api";
import type { Tenant } from "src/app/models";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import type { ContextMenuActionEvent, ContextMenuItem } from "src/app/shared/generic-table/generic-table.component";

/**
 * TenantsComponent is the controller for the table that lists Tenants.
 */
@Component({
	selector: "tp-tenants",
	styleUrls: ["./tenants.component.scss"],
	templateUrl: "./tenants.component.html"
})
export class TenantsComponent implements OnInit, OnDestroy {

	private tenantMap: Record<number, Tenant> = {};

	public searchText = "";
	public searchSubject = new BehaviorSubject("");

	public tenants: Array<Tenant> = [{
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
			valueGetter: (params: ValueGetterParams): string => {
				if (params.data.tenantId === null) {
					return "";
				}
				return `${this.tenantMap[params.data.parentId].name} (#${params.data.parentId})`;
			}
		}
	];

	public readonly contextMenuItems: ContextMenuItem<Readonly<Tenant>>[] = [
		{
			action: "viewDetails",
			name: "View Details"
		},
		{
			action: "openInNewTab",
			name: "Open in New Tab"
		}
	];

	public loading = true;
	private subscription!: Subscription;

	constructor(private readonly userService: UserService, private readonly auth: CurrentUserService) {}

	/**
	 * Angular lifecycle hook; fetches API data.
	 */
	public async ngOnInit(): Promise<void> {
		this.tenants = await this.userService.getTenants();
		this.tenantMap = Object.fromEntries((this.tenants).map(t => [t.id, t]));
		this.subscription = this.auth.userChanged.subscribe(
			() => {
				if (this.auth.hasPermission("USER:READ")) {
					this.contextMenuItems.push({
						action: "viewUsers",
						multiRow: true,
						name: "View Users"
					});
				}
				if (this.auth.hasPermission("TENANT:UPDATE")) {
					this.contextMenuItems.push({
						action: "disable",
						disabled: (ts): boolean => ts.some(t=>t.name === "root" || t.id === this.auth.currentUser?.tenantId),
						multiRow: true,
						name: "Disable"
					});
				}
			}
		);
		this.loading = false;
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
	 public handleContextMenu(a: ContextMenuActionEvent<Readonly<Tenant>>): void {
		console.log("action:", a);
	}

	/**
	 * Updates the "search" query parameter in the URL every time the search
	 * text input changes.
	 */
	public updateURL(): void {
		this.searchSubject.next(this.searchText);
	}
}
