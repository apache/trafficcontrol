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
import { DecisionDialogComponent, DecisionDialogData } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
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
	public cacheGroup: ResponseCacheGroup = {
		fallbackToClosest: true,
		fallbacks: [],
		id: -1,
		lastUpdated: new Date(),
		latitude: 0,
		localizationMethods: [],
		longitude: 0,
		name: "",
		parentCachegroupId: null,
		parentCachegroupName: null,
		secondaryParentCachegroupId: null,
		secondaryParentCachegroupName: null,
		shortName: "",
		typeId: -1,
		typeName: ""
	};

	public cacheGroups: Array<ResponseCacheGroup> = [];

	constructor(
		private readonly route: ActivatedRoute,
		private readonly api: CacheGroupService,
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
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			throw new Error(`route parameter 'id' was non-number: ${ID}`);
		}

		const cacheGroups = await this.api.getCacheGroups();
		const idx = cacheGroups.findIndex(c => c.id === numID);
		if (idx < 0) {
			throw new Error(`no such Cache Group: #${ID}`);
		}
		this.cacheGroup = cacheGroups.splice(idx, 1)[0];
		this.cacheGroups = cacheGroups;
		this.header.headerTitle.next(`Cache Group: ${this.cacheGroup.name}`);
	}

	/**
	 * Gets all Cache Groups eligible to be the parent of this Cache Group.
	 *
	 * @returns Every Cache Group except this one and its secondary parent (if
	 * it has one).
	 */
	public parentCacheGroups(): Array<ResponseCacheGroup> {
		return this.cacheGroups.filter(cg => cg.id !== this.cacheGroup.secondaryParentCachegroupId);
	}

	/**
	 * Gets all Cache Groups eligible to be the secondary parent of this Cache
	 * Group.
	 *
	 * @returns Every Cache Group except this one and its primary parent (if it
	 * has one).
	 */
	public secondaryParentCacheGroups(): Array<ResponseCacheGroup> {
		return this.cacheGroups.filter(cg => cg.id !== this.cacheGroup.parentCachegroupId);
	}

	/**
	 * Deletes the Cache Group.
	 */
	public async delete(): Promise<void> {
		if (this.new) {
			console.error("Unable to delete new Cache Group");
			return;
		}
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(
			DecisionDialogComponent,
			{
				data: {
					message: `Are you sure you want to delete Cache Group ${this.cacheGroup.name} (#${this.cacheGroup.id})?`,
					title: "Confirm Delete"
				}
			}
		);
		ref.afterClosed().subscribe(result => {
			if (result) {
				this.api.deleteCacheGroup(this.cacheGroup);
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
		if (this.new) {
			this.cacheGroup = await this.api.createCacheGroup(this.cacheGroup);
			this.new = false;
		} else {
			this.cacheGroup = await this.api.updateCacheGroup(this.cacheGroup);
		}
	}
}
