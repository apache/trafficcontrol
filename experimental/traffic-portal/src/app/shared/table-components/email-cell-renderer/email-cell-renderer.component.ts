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
import { DomSanitizer, type SafeUrl } from "@angular/platform-browser";
import type { ICellRendererParams } from "ag-grid-community";

/**
 * A table cell renderer for cells containing email addresses - it link-ifies
 * the email into a `mailto` link.
 */
@Component({
	selector: "tp-email-cell-renderer",
	styleUrls: ["./email-cell-renderer.component.scss"],
	templateUrl: "./email-cell-renderer.component.html"
})
export class EmailCellRendererComponent {

	/** The email address which the mailto link will point. */
	public get value(): string {
		return this.val;
	}

	/** The raw value for the email address. */
	private val = "";

	/** The mailto URL to use. */
	public get href(): SafeUrl {
		const url = `mailto:${this.val}`;
		return this.sanitizer.bypassSecurityTrustUrl(url);
	}

	constructor(private readonly sanitizer: DomSanitizer) {
	}

	/**
	 * Called when the value changes - I don't think this will ever happen.
	 *
	 * @param params The AG-Grid cell-rendering parameters (refer to AG-Grid docs).
	 * @returns 'true' if the component could be refreshed (always for this simple component) or 'false' if the component must be recreated.
	 */
	public refresh(params: ICellRendererParams): true {
		this.val = params.value;
		return true;
	}

	/**
	 * Called after ag-grid is initalized.
	 *
	 * @param params The AG-Grid cell-rendering parameters (refer to AG-Grid docs).
	 */
	public agInit(params: ICellRendererParams): void {
		this.val = params.value;
	}
}
