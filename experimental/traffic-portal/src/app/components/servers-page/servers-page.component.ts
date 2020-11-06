import { Component, OnInit } from "@angular/core";
import { FormControl } from "@angular/forms";
import { Router, ActivatedRoute } from "@angular/router";

import { first } from "rxjs/operators";

import { Server } from "../../models";
import { ServerService } from "../../services/api";
import { orderBy, fuzzyScore } from "../../utils";

/**
 * ServersPageComponent is the controller for the /servers page
 */
@Component({
	selector: "servers-page",
	styleUrls: ["./servers-page.component.scss"],
	templateUrl: "./servers-page.component.html",
})
export class ServersPageComponent implements OnInit {
	/** control for the fuzzy search box */
	public fuzzControl: FormControl;
	/** The list of all servers */
	public servers: Server[];
	/** All servers that match the current filter */
	public filteredServers: Server[];

	constructor (private readonly router: Router, private readonly route: ActivatedRoute, private readonly api: ServerService) { }

	/** initializes form controls and data subscriptions */
	public ngOnInit (): void {
		const searchParam = this.route.snapshot.queryParamMap.get("search");
		this.fuzzControl = new FormControl(searchParam || "");
		this.api.getServers().pipe(first()).subscribe(
			(r: Server[]) => {
				this.servers = orderBy(r, "hostName") as Server[];
				this.filteredServers = Array.from(this.servers);
				this.sort();
			}
		);
	}

	/** Updates the page URL to match the current search filter */
	public updateURL(e: Event): void {
		e.stopPropagation();
		this.sort();
		if (this.fuzzControl.value === "") {
			this.router.navigate([], {replaceUrl: true, queryParams: null});
		} else if (this.fuzzControl.value) {
			this.router.navigate([], {replaceUrl: true, queryParams: {search: this.fuzzControl.value}});
		}
	}

	/** provides the server's ID for tracking in the template */
	public tracker = (_, s: Server) =>s.id;

	/** sorts the filtered servers according to their "fuzzy score" */
	public sort(): void {
		this.filteredServers = this.servers.map(
			x => [x, fuzzyScore(x.hostName, this.fuzzControl.value)]
		).filter(x=>x[1]!==Infinity).sort(
			(a, b) => {
				if (a[1] > b[1]) {
					return 1;
				}
				if (a[1] < b[1]) {
					return -1;
				}
				return 0;
			}
		).map(x=>x[0]) as Server[];
	}
}
