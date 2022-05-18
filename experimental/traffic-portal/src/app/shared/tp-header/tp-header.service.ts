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
import {ReplaySubject} from "rxjs";

/**
 *
 */
@Injectable({
	providedIn: "root"
})
export class TpHeaderService {
	private readonly headerTitle: ReplaySubject<string>;
	private readonly headerHidden: ReplaySubject<boolean>;

	constructor() {
		this.headerTitle = new ReplaySubject(1);
		this.headerHidden = new ReplaySubject(1);
		this.headerHidden.next(true);
	}

	/**
	 * Gets the header title.
	 *
	 * @returns Subject containing the header title.
	 */
	public getTitle(): ReplaySubject<string> {
		return this.headerTitle;
	}

	/**
	 * Gets the header hidden state.
	 *
	 * @returns Subject containing header visibility state.
	 */
	public getHidden(): ReplaySubject<boolean> {
		return this.headerHidden;
	}

	/**
	 * Sets the title of the header
	 *
	 * @param title Title to use
	 */
	public setTitle(title: string): void {
		this.headerTitle.next(title);
	}

	/**
	 * Sets whether or nto to hide the header
	 *
	 * @param hidden Header visibility state
	 */
	public setHidden(hidden: boolean): void {
		this.headerHidden.next(hidden);
	}
}
