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
import { ResponseParameter } from "trafficops-types";

import { ParameterService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
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

	constructor(private readonly route: ActivatedRoute, private readonly parameterService: ParameterService,
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
			this.navSvc.headerTitle.next("New Parameter");
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
			console.error("route parameter 'id' was non-number: ", ID);
			return;
		}

		this.parameter = await this.parameterService.getParameters(numID);
		this.navSvc.headerTitle.next(`Parameter: ${this.parameter.name}`);
	}

	/**
	 * Deletes the current parameter.
	 */
	public async deleteParameter(): Promise<void> {
		if (this.new) {
			console.error("Unable to delete new parameter");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete parameter ${this.parameter.name}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.parameterService.deleteParameter(this.parameter.id);
				this.location.back();
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
			this.parameter = await this.parameterService.createParameter(this.parameter);
			this.new = false;
		} else {
			this.parameter = await this.parameterService.updateParameter(this.parameter);
		}
	}

}
