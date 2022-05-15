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
import type { ValueGetterParams } from "ag-grid-community";
import { BehaviorSubject } from "rxjs";

import { UserService } from "src/app/api";
import type { User } from "src/app/models";
import type { ContextMenuActionEvent, ContextMenuItem } from "src/app/shared/generic-table/generic-table.component";
import { orderBy } from "src/app/utils";

/**
 * UsersComponent is the controller for the "users" page.
 */
@Component({
	selector: "tp-users",
	styleUrls: ["./users.component.scss"],
	templateUrl: "./users.component.html"
})
export class UsersComponent implements OnInit {

	/** All (visible) users. */
	public users = new Array<User>();

	/** Emits changes to the fuzzy search text. */
	public fuzzySubject = new BehaviorSubject("");

	/** The current search text. */
	public searchText = "";

	/** Whether user data is still loading. */
	public loading = true;

	/**
	 * A map of Role IDs to their names, since the API doesn't provide Role
	 * names on user objects in responses.
	 */
	public roles = new Map<number, string>();

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "addressLine1",
			headerName: "Address Line 1",
			hide: true
		},
		{
			field: "addressLine2",
			headerName: "Address Line 2",
			hide: true
		},
		{
			field: "city",
			headerName: "City",
			hide: true
		},
		{
			field: "company",
			headerName: "Company",
			hide: false
		},
		{
			field: "country",
			headerName: "Country",
			hide: true
		},
		{
			cellRenderer: "emailCellRenderer",
			field: "email",
			headerName: "Email Address",
			hide: false
		},
		{
			field: "fullName",
			headerName: "Full Name",
			hide: false
		},
		{
			field: "gid",
			filter: "agNumberColumnFilter",
			headerName: "GID",
			hide: true,
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
			field: "localUser",
			filter: "tpBooleanFilter",
			headerName: "Local User",
			hide: true,
		},
		{
			field: "newUser",
			filter: "tpBooleanFilter",
			headerName: "New User",
			hide: true
		},
		{
			cellRenderer: "phoneNumberCellRenderer",
			field: "phoneNumber",
			headerName: "Phone Number",
			hide: true
		},
		{
			field: "postalCode",
			headerName: "Postal Code",
			hide: true
		},
		{
			field: "role",
			headerName: "Role",
			hide: false,
			valueGetter: (params: ValueGetterParams): string => this.roleDisplayString(params.data.role)
		},
		{
			field: "stateOrProvince",
			headerName: "State/Province",
			hide: true,
		},
		{
			field: "tenant",
			headerName: "Tenant",
			hide: false,
			valueGetter: (params: ValueGetterParams): string => `${params.data.tenant} (#${params.data.tenantId})`
		},
		{
			field: "uid",
			filter: "agNumberColumnFilter",
			headerName: "UID",
			hide: true
		},
		{
			field: "username",
			headerName: "Username",
			hide: false
		}
	];

	/** Definitions for the context menu items (which act on user data). */
	public contextMenuItems: Array<ContextMenuItem<User>> = [
		{
			action: "viewDetails",
			name: "View User Details"
		}
	];

	constructor(private readonly api: UserService) {
	}

	/**
	 * Initializes data like a map of role ids to their names.
	 */
	public async ngOnInit(): Promise<void> {
		this.roles = new Map((await this.api.getRoles()).map(r => [r.id, r.name]));
		this.users = orderBy(await this.api.getUsers(), "fullName");
		this.loading = false;
	}

	/**
	 * Gets a string suitable for displaying to the user for a given Role ID.
	 *
	 * @param role The ID of the Role being displayed.
	 * @returns A human-readable identifier for the Role, in the form `{{name}} (#{{$ID}})`.
	 */
	public roleDisplayString(role: number): string {
		const roleName = this.roles.get(role);
		if (!roleName) {
			throw new Error(`unknown Role: #${role}`);
		}
		return `${roleName} (#${role})`;
	}

	/**
	 * Updates the "search" query parameter in the URL every time the search
	 * text input changes.
	 */
	public updateURL(): void {
		this.fuzzySubject.next(this.searchText);
	}

	/**
	 * Handles the selection of a context menu item on the table.
	 *
	 * @param e The clicked action.
	 */
	public handleContextMenu(e: ContextMenuActionEvent<User>): void {
		switch (e.action) {
			case "viewDetails":
				console.log("viewing user details not implemented");
				break;
			default:
				throw new Error(`unknown context menu item clicked: ${e.action}`);
		}
	}
}
