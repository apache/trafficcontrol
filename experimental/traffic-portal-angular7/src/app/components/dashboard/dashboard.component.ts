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
import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { first } from 'rxjs/operators';

import { APIService } from '../../services';
import { DeliveryService } from '../../models';
import { orderBy } from '../../utils';

@Component({
	selector: 'dash',
	templateUrl: './dashboard.component.html',
	styleUrls: ['./dashboard.component.scss']
})
/**
 * Controller for the dashboard. Doesn't do much yet.
*/
export class DashboardComponent implements OnInit {
	deliveryServices: DeliveryService[];
	loading = true;

	now: Date;
	today: Date;

	// Fuzzy search control
	fuzzControl = new FormControl('', {updateOn: 'change'});

	constructor (private readonly api: APIService, private readonly route: ActivatedRoute, private readonly router: Router) {
		this.now = new Date();
		this.now.setUTCMilliseconds(0);
		this.today = new Date(this.now.getFullYear(), this.now.getMonth(), this.now.getDate());
	}

	ngOnInit () {
		this.api.getDeliveryServices().pipe(first()).subscribe(
			(r: DeliveryService[]) => {
				this.deliveryServices = orderBy(r, 'displayName') as DeliveryService[];
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
	}

	updateURL (e: Event) {
		e.preventDefault();
		if (this.fuzzControl.value === '') {
			this.router.navigate([], {replaceUrl: true, queryParams: null});
		} else if (this.fuzzControl.value) {
			this.router.navigate([], {replaceUrl: true, queryParams: {search: this.fuzzControl.value}});
		}
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

}
