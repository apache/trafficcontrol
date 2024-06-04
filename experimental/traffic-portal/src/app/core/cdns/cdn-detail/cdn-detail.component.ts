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
import { ResponseCDN } from "trafficops-types";

import { CDNService } from "src/app/api";
import {
	DecisionDialogComponent,
	DecisionDialogData,
} from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * CDNDetailComponent is the controller for a CDN's "detail" page.
 */
@Component({
	selector: "tp-cdn-detail",
	styleUrls: ["./cdn-detail.component.scss"],
	templateUrl: "./cdn-detail.component.html",
})
export class CDNDetailComponent implements OnInit {
	public new = false;
	public cdn: ResponseCDN = {
		dnssecEnabled: false,
		domainName: "",
		id: -1,
		lastUpdated: new Date(),
		name: "",
	};
	public showErrors = false;
	public cdns: Array<ResponseCDN> = [];

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly api: CDNService,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) {
	}

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			this.log.error("missing required route parameter 'id'");
			return;
		}

		const cdnsPromise = this.api.getCDNs().then(cdns => this.cdns = cdns);
		if (ID === "new") {
			this.new = true;
			this.setTitle();
			await cdnsPromise;
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			throw new Error(`route parameter 'id' was non-number: ${ID}`);
		}
		await cdnsPromise;
		const index = this.cdns.findIndex(c => c.id === numID);
		if (index < 0) {
			this.log.error(`no such CDN: #${ID}`);
			return;
		}
		this.cdn = this.cdns.splice(index, 1)[0];
		this.setTitle();
	}

	/**
	 * Sets the title of the page to either "new" or the name of the displayed
	 * CDN, depending on the value of
	 * {@link CDNDetailComponent.new}.
	 */
	private setTitle(): void {
		const title = this.new ? "New CDN" : `CDN: ${this.cdn.name}`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * Deletes the CDN.
	 */
	public async delete(): Promise<void> {
		if (this.new) {
			this.log.error("Unable to delete new CDN");
			return;
		}
		const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(
			DecisionDialogComponent,
			{
				data: {
					message: `Are you sure you want to delete CDN ${this.cdn.name} (#${this.cdn.id})?`,
					title: "Confirm Delete"
				}
			}
		);
		ref.afterClosed().subscribe(result => {
			if (result) {
				this.api.deleteCDN(this.cdn);
				this.router.navigate(["core/cdns"]);
			}
		});
	}

	/**
	 * Submits new/updated CDN.
	 *
	 * @param e HTML form submission event.
	 */
	public async submit(e: Event): Promise<void> {
		e.preventDefault();
		e.stopPropagation();
		this.showErrors = true;
		if (this.new) {
			this.cdn = await this.api.createCDN(this.cdn);
			this.new = false;
			await this.router.navigate(["core/cdns", this.cdn.id]);
		} else {
			this.cdn = await this.api.updateCDN(this.cdn);
		}
		this.setTitle();
	}
}
