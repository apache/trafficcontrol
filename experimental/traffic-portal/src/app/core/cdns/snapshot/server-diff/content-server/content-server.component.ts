import { Component, Input } from "@angular/core";
import { SnapshotContentServer } from "trafficops-types";

/**
 * Controller for displaying a simple Snapshot content server.
 */
@Component({
	selector: "tp-content-server[server]",
	styleUrls: ["./content-server.component.scss"],
	templateUrl: "./content-server.component.html",
})
export class ContentServerComponent {

	@Input() public server!: SnapshotContentServer;
	@Input() public kind: "unchanged" | "new" | "deleted" = "unchanged";

	/**
	 * Yields iterable tuples for the described server's `deliveryServices`. If
	 * the server isn't an edge (or hasn't been set yet by Angular), this simply
	 * returns an empty collection.
	 *
	 * @returns Tuples where the first element is the DS's XMLID and the second
	 * is the DS value.
	 */
	public deliveryServices(): Iterable<[string, Array<string>]> {
		if (!this.server || this.server.type === "MID") {
			return [];
		}
		return Object.entries(this.server.deliveryServices);
	}

}
