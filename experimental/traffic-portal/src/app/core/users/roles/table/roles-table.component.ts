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
import { BehaviorSubject } from "rxjs";
import type { ResponseRole } from "trafficops-types";

import { UserService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import type { ContextMenuActionEvent, ContextMenuItem } from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
/**
 * RolesTableComponent is the controller for the "Roles" table.
 */
@Component({
	selector: "tp-roles",
	styleUrls: ["./roles-table.component.scss"],
	templateUrl: "./roles-table.component.html"
})
export class RolesTableComponent implements OnInit {
	/** List of roles */
	public roles: Promise<Array<ResponseRole>>;
	constructor(private readonly route: ActivatedRoute, private readonly headerSvc: NavigationService,
		private readonly api: UserService, public readonly auth: CurrentUserService) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.roles = this.api.getRoles();
		this.headerSvc.headerTitle.next("Roles");
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
	}

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "name",
			headerName: "Name"
		},
		{
			field: "description",
			headerName: "Description",
		},
		{
			field: "lastUpdated",
			headerName: "Last Updated"
		}
	];

	/** Definitions for the context menu items (which act on augmented roles data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseRole>> = [
		{
			href: (selectedRow: ResponseRole): string => `${selectedRow.name}`,
			name: "View Role"
		},
		{
			href: (selectedRow: ResponseRole): string => `${selectedRow.name}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			href: (selectedRow: ResponseRole): string => `/core/users?role=${selectedRow.name}`,
			name: "View Users"
		},
	];

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
	 * @param a The action selected from the context menu.
	 */
	public handleContextMenu(a: ContextMenuActionEvent<Readonly<ResponseRole>>): void {
		console.log("action:", a);
	}
}
