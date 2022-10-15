import { Component, Input } from "@angular/core";
import { SnapshotContentRouter } from "trafficops-types";

/**
 * Controller for a component that's used to display a new, deleted, or
 * unchanged Traffic Router within a CDN Snapshot diff.
 */
@Component({
	selector: "tp-content-router",
	styleUrls: ["./content-router.component.scss"],
	templateUrl: "./content-router.component.html",
})
export class ContentRouterComponent {

	@Input() public router!: SnapshotContentRouter;
	@Input() public kind: "unchanged" | "new" | "deleted" = "unchanged";

}
