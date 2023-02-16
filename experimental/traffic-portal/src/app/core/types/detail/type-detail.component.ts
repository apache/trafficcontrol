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
import { TypeFromResponse } from "trafficops-types";

import { TypeService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * TypeDetailsComponent is the controller for the type add/edit form.
 */
@Component({
	selector: "tp-types-detail",
	styleUrls: ["./type-detail.component.scss"],
	templateUrl: "./type-detail.component.html"
})
export class TypeDetailComponent implements OnInit {
	public new = false;
	public type!: TypeFromResponse;

	constructor(private readonly route: ActivatedRoute, private readonly typeService: TypeService,
		private readonly location: Location, private readonly dialog: MatDialog, private readonly navSvc: NavigationService) { }

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
			this.navSvc.headerTitle.next("New Type");
			this.new = true;
			this.type = {
				description: "",
				id: -1,
				lastUpdated: new Date(),
				name: "",
				useInTable: "server"
			};
			return;
		}

		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			console.error("route parameter 'id' was non-number: ", ID);
			return;
		}

		this.type = await this.typeService.getTypes(numID);
		this.navSvc.headerTitle.next(`Type: ${this.type.name}`);
	}

	/**
	 * Deletes the current type.
	 */
	public async deleteType(): Promise<void> {
		if (this.new) {
			console.error("Unable to delete new type");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete type ${this.type.name}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.typeService.deleteType(this.type.id);
				this.location.back();
			}
		});
	}

	/**
	 * Submits new/updated type.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if(this.new) {
			this.type = await this.typeService.createType(this.type);
			this.new = false;
		} else {
			this.type = await this.typeService.updateType(this.type);
		}
	}

}
