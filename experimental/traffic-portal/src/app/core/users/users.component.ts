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
import { animate, style, transition, trigger } from "@angular/animations";
import { Component, type OnInit } from "@angular/core";
import type { ValueGetterParams } from "ag-grid-community";
import { BehaviorSubject } from "rxjs";
import { GetResponseUser } from "trafficops-types";

import { UserService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import type { ContextMenuItem } from "src/app/shared/generic-table/generic-table.component";
import {TpHeaderService} from "src/app/shared/tp-header/tp-header.service";
import { orderBy } from "src/app/utils";

const ANIMATION_DURATION = "150ms";

/**
 * UsersComponent is the controller for the "users" page.
 */
@Component({
	animations: [
		trigger("fadeInOut", [
			transition(
				":enter",
				[
					style({
						opacity: 0,
					}),
					animate(`${ANIMATION_DURATION} ease`, style({
						opacity: 1,
					}))
				]
			),
			transition(
				":leave",
				[
					style({
						opacity: 1,
					}),
					animate(`${ANIMATION_DURATION} ease`, style({
						opacity: 0,
					}))
				]
			)
		]),
		trigger("slideInOut", [
			transition(
				":enter",
				[
					style({
						transform: "translateY(60px)"
					}),
					animate(`${ANIMATION_DURATION} ease`, style({
						transform: "translateY(0)"
					}))
				]
			),
			transition(
				":leave",
				[
					style({
						transform: "translateY(0)"
					}),
					animate(
						`${ANIMATION_DURATION} ease`,
						style({
							transform: "translateY(60px)"
						}),
					)
				]
			)
		])
	],
	selector: "tp-users",
	styleUrls: ["./users.component.scss"],
	templateUrl: "./users.component.html"
})
export class UsersComponent implements OnInit {

	/** All (visible) users. */
	public users = new Array<GetResponseUser>();

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
			headerName: "Name",
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
	public contextMenuItems: Array<ContextMenuItem<GetResponseUser>> = [
		{
			disabled: (us: GetResponseUser | Array<GetResponseUser>): boolean => Array.isArray(us),
			href: (u: GetResponseUser): string => `/core/users/${u.id}`,
			name: "View User Details"
		},
		{
			disabled: (us: GetResponseUser | Array<GetResponseUser>): boolean => Array.isArray(us),
			href: (u: GetResponseUser): string => `/core/users/${u.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			disabled: (us: GetResponseUser | Array<GetResponseUser>): boolean => Array.isArray(us),
			href: (u: GetResponseUser): string => `/core/change-logs?search=${u.username}`,
			name: "View User Changelogs"
		}
	];

	public menuIsOpen = false;

	constructor(
		private readonly api: UserService,
		private readonly headerSvc: TpHeaderService,
		private readonly currentUserService: CurrentUserService
	) {
	}

	/**
	 * Initializes data like a map of role ids to their names.
	 */
	public async ngOnInit(): Promise<void> {
		this.roles = new Map((await this.api.getRoles()).map(r => [r.id, r.name]));
		this.users = orderBy(await this.api.getUsers(), "fullName");
		this.loading = false;
		this.headerSvc.headerTitle.next("Users");
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
	 * Toggles the state of the menu (this doesn't control the menu itself, just
	 * styling).
	 *
	 * @param closed If `"closed"`, the menu is closed, if `"opened"` it's
	 * opened.
	 */
	public toggleMenu(closed: "opened" |"closed"): void {
		if (closed === "closed") {
			this.menuIsOpen = false;
		} else {
			this.menuIsOpen = true;
		}
	}

	/**
	 * Checks if the user has permissions to create users.
	 *
	 * @returns Whether the currently authenticated user has the "USER:CREATE"
	 * Permission.
	 */
	public canCreateUsers(): boolean {
		return this.currentUserService.hasPermission("USER:CREATE");
	}
}
