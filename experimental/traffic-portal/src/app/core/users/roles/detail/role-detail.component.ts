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
import { ResponseRole } from "trafficops-types";

import { UserService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * AsnDetailComponent is the controller for the ASN add/edit form.
 */
@Component({
	selector: "tp-role-detail",
	styleUrls: ["./role-detail.component.scss"],
	templateUrl: "./role-detail.component.html"
})
export class RoleDetailComponent implements OnInit {
	public new = false;
	public asn!: ResponseRole;
	constructor(private readonly route: ActivatedRoute, private readonly location: Location,
				private readonly dialog: MatDialog, private readonly header: NavigationService) {
	}

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		const role = this.route.snapshot.paramMap.get("name");
		if (role === null) {
			console.error("missing required route parameter 'name'");
			return;
		}

		if (role === "new") {
			this.header.headerTitle.next("New Role");
			this.new = true;
			this.role = {
				description: "Read Only",
				lastUpdated: new Date(),
				name: "test",
				permissions: []
			};
			return;
		}

		this.role = await this.UserService.getRoles(role);
		this.header.headerTitle.next(`Role: ${this.role.name}`);
	}

	/**
	 * Deletes the current ASN.
	 */
	public async deleteRole(): Promise<void> {
		if (this.new) {
			console.error("Unable to delete new role");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete role ${this.role.name} with description ${this.role.description}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.UserService.deleteRole(this.role.name);
				this.location.back();
			}
		});
	}

	/**
	 * Submits new/updated ASN.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if(this.new) {
			this.asn = await this.UserService.createRole(this.role);
			this.new = false;
		} else {
			this.asn = await this.UserService.updateRoleN(this.Role);
		}
	}

}
