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
import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormControl } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { Subscription } from 'rxjs';
import { first } from 'rxjs/operators';

import { APIService, AuthenticationService } from '../../services';
import { DeliveryService } from '../../models';
import { orderBy, fuzzyScore } from '../../utils';

@Component({
	selector: 'dash',
	templateUrl: './dashboard.component.html',
	styleUrls: ['./dashboard.component.scss']
})
/**
 * Controller for the dashboard. Doesn't do much yet.
*/
export class DashboardComponent implements OnInit, OnDestroy {
	deliveryServices: DeliveryService[];
	filteredDSes: DeliveryService[];
	loading = true;

	private capabilitiesSubscription: Subscription;
	canCreateDeliveryServices = false;

	now: Date;
	today: Date;

	// Fuzzy search control
	fuzzControl = new FormControl('', {updateOn: 'change'});

	constructor (private readonly api: APIService, private readonly route: ActivatedRoute, private readonly router: Router, private readonly auth: AuthenticationService) {
		this.now = new Date();
		this.now.setUTCMilliseconds(0);
		this.today = new Date(this.now.getFullYear(), this.now.getMonth(), this.now.getDate());
	}

	ngOnInit () {
		this.api.getDeliveryServices().pipe(first()).subscribe(
			(r: DeliveryService[]) => {
				this.deliveryServices = orderBy(r, 'displayName') as DeliveryService[];
				this.filteredDSes = Array.from(this.deliveryServices);
				this.loading = false;
			}
		);

		this.route.queryParamMap.pipe(first()).subscribe(
			m => {
				if (m.has('search')) {
					this.fuzzControl.setValue(decodeURIComponent(m.get('search')));
				}
			}
		);

		this.capabilitiesSubscription = this.auth.currentUserCapabilities.subscribe(
			v => {
				this.canCreateDeliveryServices = v.has("ds-create");
			}
		)
	}

	updateURL (e: Event) {
		e.preventDefault();
		this.sort();
		if (this.fuzzControl.value === '') {
			this.router.navigate([], {replaceUrl: true, queryParams: null});
		} else if (this.fuzzControl.value) {
			this.router.navigate([], {replaceUrl: true, queryParams: {search: this.fuzzControl.value}});
		}
	}

	sort () {
		this.filteredDSes = this.deliveryServices.map(x=>[x, fuzzyScore(x.displayName, this.fuzzControl.value)]).filter(x=>x[1] !== Infinity).sort(
			(a, b) => {
				if (a[1] > b[1]) {
					return 1;
				}
				if (a[1] < b[1]) {
					return -1;
				}
				return 0;
			}
		).map(x=>x[0]) as DeliveryService[];
	}

	/**
	 * Checks if a Delivery Service matches a fuzzy search term
	 * @param ds The Delivery Service being checked
	 * @returns `true` if `ds` matches the fuzzy search term, `false` otherwise
	*/
	fuzzy (ds: DeliveryService): boolean {
		if (!this.fuzzControl.value) {
			return true;
		}
		const testVal = ds.displayName.toLocaleLowerCase();
		let n = -1;
		for (const l of this.fuzzControl.value.toLocaleLowerCase()) {
			/* tslint:disable */
			if (!~(n = testVal.indexOf(l, n + 1))) {
			/* tslint:enable */
				return false;
			}
		}
		return true;
	}

	tracker (unused_item: number, d: DeliveryService) {
		return d.id;
	}

	ngOnDestroy () {
		this.capabilitiesSubscription.unsubscribe();
	}

}
