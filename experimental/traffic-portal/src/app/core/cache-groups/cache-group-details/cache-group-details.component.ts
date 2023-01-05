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

import { Location } from "@angular/common";
import { Component, OnInit } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import type { ResponseCacheGroup } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { TpHeaderService } from "src/app/shared/tp-header/tp-header.service";

/**
 * The controller for the form page for creating/updating a Cache Group.
 */
@Component({
	selector: "tp-cache-group-details",
	styleUrls: ["./cache-group-details.component.scss"],
	templateUrl: "./cache-group-details.component.html",
})
export class CacheGroupDetailsComponent implements OnInit {
	public new = false;
	public cacheGroup!: ResponseCacheGroup;

	constructor(
		private readonly route: ActivatedRoute,
		private readonly cacheGroupService: CacheGroupService,
		private readonly location: Location,
		private readonly dialog: MatDialog,
		private readonly header: TpHeaderService
	) { }

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			console.error("missing required route parameter 'id'");
			return;
		}

		if (ID === "new") {
			this.header.headerTitle.next("New Division");
			this.new = true;
			this.cacheGroup = {
				fallbackToClosest: true,
				fallbacks: [],
				id: -1,
				lastUpdated: new Date(),
				latitude: 0,
				localizationMethods: [],
				longitude: 0,
				name: "",
				parentCacheGroupId: null,
				parentCacheGroupName: null,
				secondaryParentCacheGroupId: null,
				secondaryParentCacheGroupName: null,
				shortName: "",
				typeId: -1,
				typeName: ""
			};
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			console.error("route parameter 'id' was non-number:", ID);
			return;
		}

		this.cacheGroup = await this.cacheGroupService.getCacheGroups(numID);
		this.header.headerTitle.next(`Cache Group: ${this.cacheGroup.name}`);
	}

	/**
	 * Deletes the Cache Group.
	 */
	public async deleteDivision(): Promise<void> {
		if (this.new) {
			console.error("Unable to delete new Cache Group");
			return;
		}
		const ref = this.dialog.open(
			DecisionDialogComponent,
			{
				data: {
					message: `Are you sure you want to delete Cache Group ${this.cacheGroup.name} (#${this.cacheGroup.id})?`,
					title: "Confirm Delete"
				}
			}
		);
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.cacheGroupService.deleteCacheGroup(this.cacheGroup);
				this.location.back();
			}
		});
	}

	/**
	 * Submits new/updated Cache Group.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		this.cacheGroup.shortName = this.cacheGroup.name;
		if(this.new) {
			this.cacheGroup = await this.cacheGroupService.createCacheGroup(this.cacheGroup);
			this.new = false;
		} else {
			this.cacheGroup = await this.cacheGroupService.updateCacheGroup(this.cacheGroup);
		}
	}
}
