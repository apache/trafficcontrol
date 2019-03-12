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
import { Component, Input } from '@angular/core';

import { first } from 'rxjs/operators';

import { APIService } from '../../services';
import { DeliveryService } from '../../models/deliveryservice';

@Component({
	selector: 'ds-card',
	templateUrl: './ds-card.component.html',
	styleUrls: ['./ds-card.component.scss']
})
export class DsCardComponent {

	@Input() deliveryService: DeliveryService;

	// Capacity measures
	available: number;
	maintenance: number;
	utilized: number;

	private loaded: boolean;

	constructor(private api: APIService) {
		this.available = 100;
		this.maintenance = 0;
		this.utilized = 0;
		this.loaded = false;
	}

	toggle(e: Event) {
		if (!this.loaded && (e.target as HTMLDetailsElement).open) {
			this.api.getDSCapacity(this.deliveryService.id).pipe(first()).subscribe(
				r => {
					if (r) {
						this.available = r.availablePercent;
						this.maintenance = r.maintenancePercent;
						this.utilized = r.utilizedPercent;
						this.loaded = true;
					}
				}
			);
		}
	}

}
