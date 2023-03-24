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
import { Component, OnInit } from "@angular/core";
import { BehaviorSubject } from "rxjs";
import { ResponseStatus } from "trafficops-types";

import { ServerService } from "src/app/api";
import { ContextMenuItem } from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * StatusesTableComponent is the controller for the statuses page - which
 * principally contains a table.
 */
@Component({
	selector: "tp-statuses-table",
	styleUrls: ["./statuses-table.component.scss"],
	templateUrl: "./statuses-table.component.html",
})
export class StatusesTableComponent implements OnInit {

	/** All of the statues which should appear in the table. */
	public statuses: ResponseStatus[] = [];

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [
		{
			field: "name",
			headerName: "Name",
			hide: false
		},
		{
			field: "description",
			headerName: "Description",
			hide: false
		}];

	/** The current search text. */
	public searchText = "";

	/** Definitions for the context menu items (which act on statuses data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseStatus>> = [
		{
			href: (u: ResponseStatus): string => `${u.id}`,
			name: "View Status Details"
		},
		{
			href: (): string => "new",
			name: "Create New Status"
		}
	];

	/** Emits changes to the fuzzy search text. */
	public fuzzySubject = new BehaviorSubject("");

	/**
	 * Constructs the component with its required injections.
	 *
	 * @param serverService The Servers API which is used to provide row data.
	 * @param navSvc Manages the header
	 */
	constructor(
		private readonly serverService: ServerService,
		private readonly navSvc: NavigationService,
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.navSvc.headerTitle.next("Statuses");
	}

	/** Initializes table data, loading it from Traffic Ops. */
	public ngOnInit(): void {
		this.getStatuses();
	}

	/** Reloads the servers table data. */
	public async getStatuses(): Promise<void> {
		this.statuses = await this.serverService.getStatuses();
	}

	/**
	 * Updates the "search" query parameter in the URL every time the search
	 * text input changes.
	 */
	public updateURL(): void {
		this.fuzzySubject.next(this.searchText);
	}
}
