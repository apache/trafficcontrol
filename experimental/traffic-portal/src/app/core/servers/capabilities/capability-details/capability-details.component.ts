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
import { ResponseServerCapability } from "trafficops-types";

import { ServerService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * Controller for the form for creating and editing Server Capabilities.
 */
@Component({
	selector: "tp-capability-details",
	styleUrls: ["./capability-details.component.scss"],
	templateUrl: "./capability-details.component.html",
})
export class CapabilityDetailsComponent implements OnInit {
	public new = false;

	public capability!: ResponseServerCapability;

	/**
	 * This caches the original name of the Capability, so that updates can be
	 * made.
	 */
	private name = "";

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly api: ServerService,
	) {}

	/**
	 * Angular lifecycle hook.
	 */
	public async ngOnInit(): Promise<void> {
		const name = this.route.snapshot.paramMap.get("name");
		if (name === null) {
			this.setHeader("New Capability");
			this.new = true;
			this.capability = {
				lastUpdated: new Date(),
				name: "",
			};
			return;
		}

		this.capability = await this.api.getCapabilities(name);
		this.name = this.capability.name;
		this.navSvc.headerTitle.next(`Capability: ${this.capability.name}`);
	}

	/**
	 * Sets the value of the header text, and caches the Capability's initial
	 * name.
	 *
	 * @param name The name of the current Capability (before editing).
	 */
	private setHeader(name: string): void {
		this.name = name;
		this.navSvc.headerTitle.next(`Capability: ${name}`);
	}

	/**
	 * Deletes the current physLocation.
	 */
	public async deleteCapability(): Promise<void> {
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {
				message: `Are you sure you want to delete the Capability '${this.capability.name}'?`,
				title: "Confirm Delete"
			}
		});
		const result = await ref.afterClosed().toPromise();
		if(result) {
			await this.api.deleteCapability(this.capability);
			await this.router.navigate(["core/capabilities"]);
		}
	}

	/**
	 * Submits new/updated physLocation.
	 *
	 * @param e HTML click event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		if(this.new) {
			this.capability = await this.api.createCapability(this.capability);
			this.new = false;
			this.setHeader("New Capability");
		} else {
			this.capability = await this.api.updateCapability(this.name, this.capability);
			this.navSvc.headerTitle.next(`Capability: ${this.capability.name}`);
		}
		this.router.navigate([`/core/capabilities/${this.capability.name}`], {replaceUrl: true});
	}
}
