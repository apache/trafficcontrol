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
import { Injectable } from "@angular/core";

import {TpHeaderComponent} from "src/app/shared/tp-header/tp-header.component";

/**
 *
 */
@Injectable({
	providedIn: "root"
})
export class TpHeaderService {
	private header?: TpHeaderComponent;

	/**
	 * Register the header component
	 *
	 * @param header The header to register
	 */
	public registerHeader(header: TpHeaderComponent): void {
		this.header = header;
	}

	/**
	 * Sets the title of the header
	 *
	 * @param title Title to use
	 */
	public setTitle(title: string): void {
		if(this.header) {
			this.header.title = title;
		}
	}

	/**
	 * Sets whether or nto to hide the header
	 *
	 * @param hidden Header visibility state
	 */
	public setHidden(hidden: boolean): void {
		if(this.header) {
			this.header.hidden = hidden;
		}
	}
}
