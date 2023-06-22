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
import { FormControl, UntypedFormControl } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Params } from "@angular/router";
import { BehaviorSubject } from "rxjs";
import { ProfileImport, ResponseProfile } from "trafficops-types";

import { ProfileService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import {
	ContextMenuActionEvent,
	ContextMenuItem,
	DoubleClickLink,
	TableTitleButton
} from "src/app/shared/generic-table/generic-table.component";
import { ImportJsonTxtComponent } from "src/app/shared/import-json-txt/import-json-txt.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * ProfileTableComponent is the controller for the profiles page - which
 * principally contains a table.
 */
@Component({
	selector: "tp-profile-table",
	styleUrls: ["./profile-table.component.scss"],
	templateUrl: "./profile-table.component.html"
})
export class ProfileTableComponent implements OnInit {
	/** All the physical locations which should appear in the table. */
	public profiles: Promise<Array<ResponseProfile>>;

	/** Definitions of the table's columns according to the ag-grid API */
	public columnDefs = [{
		field: "cdnName",
		headerName: "CDN"
	}, {
		field: "description",
		headerName: "Description",
	}, {
		field: "id",
		headerName: "ID",
		hide: true
	}, {
		field: "lastUpdated",
		headerName: "Last Updated",
		hide: true
	}, {
		field: "name",
		headerName: "Name"
	}, {
		field: "routingDisabled",
		headerName: "Routing Disabled"
	}, {
		field: "type",
		headerName: "Type"
	}];

	public titleBtns: Array<TableTitleButton> = [
		{
			action: "import",
			text: "Import Profile",
		}
	];

	/** Definitions for the context menu items (which act on augmented cache-group data). */
	public contextMenuItems: Array<ContextMenuItem<ResponseProfile>> = [
		{
			href: (profile: ResponseProfile): string => `${profile.id}`,
			name: "Open in New Tab",
			newTab: true
		},
		{
			href: (type: ResponseProfile): string => `${type.id}`,
			name: "Edit"
		},
		{
			action: "delete",
			multiRow: false,
			name: "Delete"
		},
		{
			action: "clone-profile",
			disabled: (): true => true,
			multiRow: false,
			name: "Clone Profile",
		},
		{
			href: (profile: ResponseProfile): string => `/api/${this.api.apiVersion}/profiles/${profile.id}/export`,
			name: "Export Profile",
			newTab: true
		},
		{
			action: "manage-parameters",
			disabled: (): true => true,
			multiRow: false,
			name: "Manage Parameters",
		},
		{
			href: "/core/servers",
			name: "View Servers",
			queryParams: (profile: ResponseProfile): Params => ({profileName: profile.name})
		}
	];

	/** Defines what the table should do when a row is double-clicked. */
	public doubleClickLink: DoubleClickLink<ResponseProfile> = {
		href: (row: ResponseProfile): string => `/core/profiles/${row.id}`
	};

	/** A subject that child components can subscribe to for access to the fuzzy search query text */
	public fuzzySubject: BehaviorSubject<string>;

	/** Form controller for the user search input. */
	public fuzzControl: UntypedFormControl = new FormControl<string>("", {nonNullable: true});

	/**
	 * Constructs the component with its required injections.
	 *
	 * @param api The Servers API which is used to provide row data.
	 * @param route A reference to the route of this view which is used to set the fuzzy search box text from the 'search' query parameter.
	 * @param router Angular router
	 * @param navSvc Manages the header
	 * @param dialog Dialog manager
	 */
	constructor(
		private readonly api: ProfileService,
		private readonly route: ActivatedRoute,
		private readonly navSvc: NavigationService,
		private readonly dialog: MatDialog,
		public readonly auth: CurrentUserService,
	) {
		this.fuzzySubject = new BehaviorSubject<string>("");
		this.profiles = this.api.getProfiles();
		this.navSvc.headerTitle.next("Profiles");
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
		const hasParameter = this.route.snapshot.queryParamMap.get("hasParameter");
		if (hasParameter == null) {
			return;
		}
		this.profiles = this.api.getProfilesByParam(+hasParameter);
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
	public async handleContextMenu(evt: ContextMenuActionEvent<ResponseProfile>): Promise<void> {
		const data = evt.data as ResponseProfile;
		switch (evt.action) {
			case "delete":
				const ref = this.dialog.open(DecisionDialogComponent, {
					data: {
						message: `Are you sure to delete Profile ${data.name} with id ${data.id}?`,
						title: "Confirm Delete"
					}
				});
				ref.afterClosed().subscribe(result => {
					if (result) {
						this.api.deleteProfile(data.id).then(async () => this.profiles = this.api.getProfiles());
					}
				});
				break;
		}
	}

	/**
	 * handles when a title button is event is emitted
	 *
	 * @param action which button was pressed
	 */
	public async handleTitleButton(action: string): Promise<void> {
		switch(action){
			case "import":
				const ref = this.dialog.open(ImportJsonTxtComponent,{
					data: { title: "Import Profile" },
					width: "70vw"
				});

				/** After submission from Import JSON dialog component */
				ref.afterClosed().subscribe( (result: ProfileImport) => {
					if (result) {
						this.api.importProfile(result).then(response => {
							if (response) {
								this.profiles = this.api.getProfiles();
							}
						});
					}
				});
				break;
		}
	}
}
