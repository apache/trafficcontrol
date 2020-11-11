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

import { Component } from "@angular/core";
import { DomSanitizer, SafeUrl } from "@angular/platform-browser";

import { ICellRendererAngularComp } from "ag-grid-angular";
import { ICellRendererParams } from "ag-grid-community";

import { AuthenticationService } from "../../../services";

/**
 * SSHCellRendererComponent is an AG-Grid cell renderer that provides ssh:// links as content.
 */
@Component({
	selector: "ssh-cell-renderer",
	styleUrls: ["./ssh-cell-renderer.component.scss"],
	templateUrl: "./ssh-cell-renderer.component.html"
})
export class SSHCellRendererComponent implements ICellRendererAngularComp {

	/** The IP address or hostname to which the SSH link will point. */
	public get value(): string {
		return this.val;
	}
	private val = "";

	private username = "";

	/** The SSH URL to use. */
	public get href(): SafeUrl {
		const url = `ssh://${this.username}@${this.value}`;
		return this.sanitizer.bypassSecurityTrustUrl(url);
	}

	constructor(private readonly auth: AuthenticationService, private readonly sanitizer: DomSanitizer) {
		this.auth.updateCurrentUser().subscribe(
			success => {
				if (success) {
					const cu = this.auth.currentUserValue;
					if (cu) {
						this.username = cu.username;
					}
				}
			}
		);
		const u = this.auth.currentUserValue;
		if (u) {
			this.username = u.username;
		}
	}

	/** Called when the value changes - I don't think this will ever happen. */
	public refresh(params: ICellRendererParams): boolean {
		this.val = params.value;
		console.log("refreshed:", params);
		return true;
	}

	/** called after ag-grid is initalized */
	public agInit(params: ICellRendererParams): void {
		console.log("has value?:", Object.prototype.hasOwnProperty.call(params, "value"));
		console.log("getval:", params.getValue());
		this.val = params.value;
	}
}
