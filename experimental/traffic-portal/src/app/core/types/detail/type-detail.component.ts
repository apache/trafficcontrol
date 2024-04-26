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
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import { TypeFromResponse } from "trafficops-types";

import { TypeService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
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

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly typeService: TypeService,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) { }

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			this.log.error("missing required route parameter 'id'");
			return;
		}

		this.new = ID === "new";

		if (this.new) {
			this.setTitle();
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
			this.log.error("route parameter 'id' was non-number: ", ID);
			return;
		}

		this.type = await this.typeService.getTypes(numID);
		this.setTitle();
	}

	/**
	 * Sets the headerTitle based on current Type state.
	 *
	 * @private
	 */
	private setTitle(): void {
		const title = this.new ? "New Type" : `Type: ${this.type.name}`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * Deletes the current type.
	 */
	public async deleteType(): Promise<void> {
		if (this.new) {
			this.log.error("Unable to delete new type");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete type ${this.type.name}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.typeService.deleteType(this.type.id);
				this.router.navigate(["core/types"]);
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
			await this.router.navigate(["core/types", this.type.id]);
		} else {
			this.type = await this.typeService.updateType(this.type);
		}
		this.setTitle();
	}

}
