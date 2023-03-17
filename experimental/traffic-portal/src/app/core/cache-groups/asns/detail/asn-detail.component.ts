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
import { ResponseASN, ResponseCacheGroup } from "trafficops-types";

import { CacheGroupService } from "src/app/api";
import { DecisionDialogComponent } from "src/app/shared/dialogs/decision-dialog/decision-dialog.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * AsnDetailComponent is the controller for the ASN add/edit form.
 */
@Component({
	selector: "tp-asn-detail",
	styleUrls: ["./asn-detail.component.scss"],
	templateUrl: "./asn-detail.component.html"
})
export class AsnDetailComponent implements OnInit {
	public new = false;
	public asn!: ResponseASN;
	public cachegroups!: Array<ResponseCacheGroup>;
	constructor(private readonly route: ActivatedRoute, private readonly cacheGroupService: CacheGroupService,
		private readonly location: Location, private readonly dialog: MatDialog,
		private readonly header: NavigationService) {
	}

	/**
	 * Angular lifecycle hook where data is initialized.
	 */
	public async ngOnInit(): Promise<void> {
		this.cachegroups = await this.cacheGroupService.getCacheGroups();
		const ID = this.route.snapshot.paramMap.get("id");
		if (ID === null) {
			console.error("missing required route parameter 'id'");
			return;
		}

		if (ID === "new") {
			this.header.headerTitle.next("New ASN");
			this.new = true;
			this.asn = {
				asn: 1,
				cachegroup: "test",
				cachegroupId: 1,
				id: 1,
				lastUpdated: new Date()
			};
			return;
		}
		const numID = parseInt(ID, 10);
		if (Number.isNaN(numID)) {
			console.error("route parameter 'id' was non-number:", ID);
			return;
		}

		this.asn = await this.cacheGroupService.getASNs(numID);
		this.header.headerTitle.next(`ASN: ${this.asn.asn}`);
	}

	/**
	 * Deletes the current ASN.
	 */
	public async deleteAsn(): Promise<void> {
		if (this.new) {
			console.error("Unable to delete new ASN");
			return;
		}
		const ref = this.dialog.open(DecisionDialogComponent, {
			data: {message: `Are you sure you want to delete ASN ${this.asn.asn} with id ${this.asn.id}`,
				title: "Confirm Delete"}
		});
		ref.afterClosed().subscribe(result => {
			if(result) {
				this.cacheGroupService.deleteASN(this.asn.id);
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
			this.asn = await this.cacheGroupService.createASN(this.asn);
			this.new = false;
		} else {
			this.asn = await this.cacheGroupService.updateASN(this.asn);
		}
	}

}
