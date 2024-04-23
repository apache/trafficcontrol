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
import { ResponseASN, ResponseCacheGroup } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * AsnDetailComponent is the controller for the ASN add/edit form.
 */
@Component({
	selector: "tp-asn-detail",
	styleUrls: ["./asn-detail.component.scss"],
	templateUrl: "./asn-detail.component.html"
})
export class ASNDetailComponent implements OnInit {
	public isNew = false;
	public asn!: ResponseASN;
	public cacheGroups =  new Array<ResponseCacheGroup>();

	constructor(
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly cacheGroupService: CacheGroupService,
		private readonly dialog: MatDialog,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) {
	}

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		this.cacheGroupService.getCacheGroups().then(
			cgs => {
				this.cacheGroups = cgs;
			}
		);

		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			this.log.error("missing required route parameter 'id'");
			return;
		}

		this.isNew = ID === "new";

		if (this.isNew) {
			this.setTitle();
			this.isNew = true;
			this.asn = {
				asn: 0,
				cachegroup: "",
				cachegroupId: 0,
				id: 0,
				lastUpdated: new Date()
			};
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			this.log.error("route parameter 'id' was non-number:", ID);
			return;
		}

		this.asn = await this.cacheGroupService.getASNs(numID);
		this.setTitle();
	}

	/**
	 * Sets the headerTitle based on current ASN state.
	 *
	 * @private
	 */
	private setTitle(): void {
		const title = this.isNew ? "New ASN" : `ASN: ${this.asn.asn}`;
		this.navSvc.headerTitle.next(title);
	}

	/**
	 * Deletes the current ASN.
	 */
	public async deleteAsn(): Promise<void> {
		if (this.isNew) {
			this.log.error("Unable to delete new ASN");
			return;
		}
		if(!this.asn.asn) {
			this.log.error("Missing ASN number");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete ASN ${this.asn.asn} with id ${this.asn.id}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.cacheGroupService.deleteASN(this.asn.id);
				this.router.navigate(["core/asns"]);
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
		if(this.isNew) {
			this.asn = await this.cacheGroupService.createASN(this.asn);
			this.isNew = false;
			await this.router.navigate(["core/asns", this.asn.id]);
		} else {
			this.asn = await this.cacheGroupService.updateASN(this.asn);
		}
		this.setTitle();
	}

}
