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
import { Component, OnInit, OnDestroy } from "@angular/core";
import { FormControl } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";

import { Subscription } from "rxjs";
import { first } from "rxjs/operators";

import { DeliveryService } from "../../models";
import { AuthenticationService } from "../../services";
import { DeliveryServiceService } from "../../services/api";
import { orderBy, fuzzyScore } from "../../utils";

/**
 * DashboardComponent is the controller for the dashboard, where a user sees all
 * of the Delivery Services in their Tenancy.
 */
@Component({
	selector: "tp-dash",
	styleUrls: ["./dashboard.component.scss"],
	templateUrl: "./dashboard.component.html"
})
export class DashboardComponent implements OnInit, OnDestroy {
	/**
	 * The set of all Delivery Services (visible to the Tenant).
	 */
	public deliveryServices: DeliveryService[] = [];

	/**
	 * The set of Delivery Services filtered according to the search box text.
	 */
	public get filteredDSes(): DeliveryService[] {
		if (!this.deliveryServices) {
			return [];
		}
		return this.deliveryServices.map(x => [x, fuzzyScore(x.displayName, this.fuzzControl.value)]).filter(x => x[1] !== Infinity).sort(
			(a, b) => {
				if (a[1] > b[1]) {
					return 1;
				}
				if (a[1] < b[1]) {
					return -1;
				}
				return 0;
			}
		).map(x => x[0]) as Array<DeliveryService>;
	}

	/** Whether or not the page is still loading. */
	public loading = true;

	/** A subscription for the Capabilities of the user. */
	private capabilitiesSubscription: Subscription | null = null;

	/**
	 * Whether or not the currently logged-in user has permission to create
	 * Delivery Services.
	 */
	public canCreateDeliveryServices = false;

	/**
	 * The date and time at which the page loaded.
	 */
	public now: Date;

	/** 00:00:00 GMT on the day the page loaded. */
	public today: Date;

	/** Fuzzy search control */
	public fuzzControl = new FormControl("", {updateOn: "change"});

	/**
	 * Constructor.
	 */
	constructor(
		private readonly dsAPI: DeliveryServiceService,
		private readonly route: ActivatedRoute,
		private readonly router: Router,
		private readonly auth: AuthenticationService
	) {
		this.now = new Date();
		this.now.setUTCMilliseconds(0);
		this.today = new Date(this.now.getFullYear(), this.now.getMonth(), this.now.getDate());
	}

	/**
	 * Runs initialization, fetching the list of (visible) Delivery Services and
	 * setting the search test from the query parameters, if applicable.
	 */
	public ngOnInit(): void {
		this.dsAPI.getDeliveryServices().pipe(first()).subscribe(
			r => {
				// these annoying typecasts are necessary because of how object property indexing works.
				// look at 'orderBy' to understand.
				this.deliveryServices = (
					orderBy((r as unknown[]) as Record<string, unknown>[], "displayName") as unknown[]
				) as DeliveryService[];
				this.loading = false;
			}
		);

		this.route.queryParamMap.pipe(first()).subscribe(
			m => {
				const search = m.get("search");
				if (search) {
					this.fuzzControl.setValue(decodeURIComponent(search));
				}
			}
		);

		this.capabilitiesSubscription = this.auth.currentUserCapabilities.subscribe(
			v => {
				this.canCreateDeliveryServices = v.has("ds-create");
			}
		);
	}

	/**
	 * Updates the page's URL to show the current search term in its 'search'
	 * query parameter.
	 *
	 * @param e The Event that triggered updating the URL.
	 */
	public updateURL(e: Event): void {
		e.preventDefault();
		if (this.fuzzControl.value === "") {
			this.router.navigate([], {queryParams: null, replaceUrl: true});
		} else if (this.fuzzControl.value) {
			this.router.navigate([], {queryParams: {search: this.fuzzControl.value}, replaceUrl: true});
		}
	}

	/**
	 * Provides a tracking identifier for each Delivery Service.
	 *
	 * @param _ unused, but specified by the Angular API.
	 * @param d The Delivery Service to identify.
	 * @returns A unique identifier for 'd' to facilitate tracking.
	 */
	public tracker(_: number, d: DeliveryService): number {
		return d.id || 0;
	}

	/**
	 * Cleans up subscriptions when the component is destroyed.
	 */
	public ngOnDestroy(): void {
		if (this.capabilitiesSubscription) {
			this.capabilitiesSubscription.unsubscribe();
		}
	}

}
