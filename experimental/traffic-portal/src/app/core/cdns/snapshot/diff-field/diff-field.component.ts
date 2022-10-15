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
