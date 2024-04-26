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
import { ResponseParameter } from "trafficops-types";

import { ProfileService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * ParameterDetailsComponent is the controller for the parameter add/edit form.
 */
@Component({
	selector: "tp-parameters-detail",
	styleUrls: ["./parameter-detail.component.scss"],
	templateUrl: "./parameter-detail.component.html"
})
export class ParameterDetailComponent implements OnInit {
	public new = false;
	public parameter!: ResponseParameter;
	public secure = [
		{ label: "true", value: true },
		{ label: "false", value: false }
	];

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly profileService: ProfileService,
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
			this.parameter = {
				configFile: "",
				id: -1,
				lastUpdated: new Date(),
				name: "",
				profiles: [],
				secure: false,
				value: "",
			};
			return;
		}

		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			this.log.error("route parameter 'id' was non-number: ", ID);
			return;
		}

		this.parameter = await this.profileService.getParameters(numID);
		this.setTitle();
	}

	/**
	 * Sets the headerTitle based on current Parameter state.
	 *
	 * @private
	 */
	private setTitle(): void {
		const title = this.new ? "New Parameter" : `Parameter: ${this.parameter.name} (${this.parameter.id})`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * Deletes the current parameter.
	 */
	public async deleteParameter(): Promise<void> {
		if (this.new) {
			this.log.error("Unable to delete new parameter");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete parameter ${this.parameter.name}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.profileService.deleteParameter(this.parameter.id);
				this.router.navigate(["core/parameters"]);
			}
		});
	}

	/**
	 * Submits new/updated parameter.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if(this.new) {
			this.parameter = await this.profileService.createParameter(this.parameter);
			this.new = false;
			await this.router.navigate(["core/parameters", this.parameter.id]);
		} else {
			this.parameter = await this.profileService.updateParameter(this.parameter);
		}
		this.setTitle();
	}
}
