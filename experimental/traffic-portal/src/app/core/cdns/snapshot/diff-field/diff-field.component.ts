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

import { Component, Input } from "@angular/core";

import type { DiffVal } from "src/app/utils/snapshot.diffing";

/**
 * A "diff field" is a component used for displaying generically the difference
 * - or lack thereof - between a simple property in the current and pending
 * Snapshots of a CDN. It can't be used to express differences in non-primitive
 * types.
 */
@Component({
	selector: "tp-diff-field[value]",
	styleUrls: ["./diff-field.component.scss"],
	templateUrl: "./diff-field.component.html",
})
export class DiffFieldComponent {

	@Input() public last = false;
	@Input() public value!: DiffVal;
	@Input() public name?: string;

}
