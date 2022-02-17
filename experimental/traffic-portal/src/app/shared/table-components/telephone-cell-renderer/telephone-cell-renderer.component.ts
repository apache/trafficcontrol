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
import { ICellRendererParams } from "ag-grid-community";

/**
 * A table cell renderer for cells containing telephone numbers - it link-ifies
 * the phone number into a `tel` link.
 */
@Component({
	selector: "tp-telephone-cell-renderer",
	styleUrls: ["./telephone-cell-renderer.component.scss"],
	templateUrl: "./telephone-cell-renderer.component.html"
})
export class TelephoneCellRendererComponent {

	/** The telephone number which the tel link will point. */
	public get value(): string {
		return this.val;
	}

	/** The raw value for the phone number. */
	private val = "";

	/** The URL to use for the link. */
	public get href(): SafeUrl {
		const url = `tel:${this.val}`;
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
