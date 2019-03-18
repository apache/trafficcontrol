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
import { Component, Input, ElementRef } from '@angular/core';

import { first } from 'rxjs/operators';

import { Chart } from 'chart.js';

import { APIService } from '../../services';
import { DeliveryService, Protocol } from '../../models/deliveryservice';

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

	// Health measured as a percent of Cache Groups that are healthy
	healthy: number;

	// Bandwidth data
	chart: Chart;
	// midBandwidth: Array<Array<any>>;
	edgeBandwidth: Array<number>;
	labels: Array<Date>;

	// Need this to access merged namespace for string conversions
	Protocol = Protocol;


	private loaded: boolean;
	public graphDataLoaded: boolean;

	constructor(private api: APIService, private elementRef: ElementRef) {
		this.available = 100;
		this.maintenance = 0;
		this.utilized = 0;
		this.loaded = false;
		this.edgeBandwidth = new Array<number>();
		this.labels = new Array<Date>();
		this.graphDataLoaded = false;
	}

	/**
	 * Triggered when the details element is opened or closed, and fetches the latest stats.
	 * @param {e} A DOM Event caused then the open/close state of a `<details>` element changes.
	 *
	 * this will only fetch health and capacity information once per page load, but will update the
	 * Bandwidth graph every time the details panel is opened. Bandwidth data is calculated using
	 * 60s intervals starting at 00:00 UTC the current day and ending at the current date/time.
	*/
	toggle(e: Event) {
		if ((e.target as HTMLDetailsElement).open) {
			if (!this.loaded) {
				this.loaded = true;
				this.api.getDSCapacity(this.deliveryService.id).pipe(first()).subscribe(
					r => {
						if (r) {
							this.available = r.availablePercent;
							this.maintenance = r.maintenancePercent;
							this.utilized = r.utilizedPercent;
						}
					}
				);
				this.api.getDSHealth(this.deliveryService.id).pipe(first()).subscribe(
					r => {
						if (r) {
							this.healthy = r.totalOnline / (r.totalOnline + r.totalOffline);
						}
					}
				);
			} else if (this.chart) {
				this.chart.destroy();
				this.graphDataLoaded = false;
			}
			const now = new Date();
			now.setUTCMilliseconds(0); // apparently `const` doesn't care about this
			const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
			this.api.getDSKBPS(this.deliveryService.xmlId, today, now, '60s').pipe(first()).subscribe(
				data => {
					for (let d of data) {
						this.labels.push(new Date(d[0]));
						this.edgeBandwidth.push(d[1]);
					}

					const canvas = this.elementRef.nativeElement.querySelector('#canvas-'+String(this.deliveryService.id));
					this.chart = new Chart(canvas, {
						type: 'line',
						data: {
							labels: this.labels,
							datasets: [
								{
									data: this.edgeBandwidth,
									borderColor: '#3cba9f',
									fill: false
								}
							]
						},
						options: {
							legend: {
								display: false
							},
							title: {
								display: true,
								text: "Today's Bandwidth"
							},
							scales: {
								xAxes: [{
									display: true,
									type: 'time',
									callback: (v, i, values) => {
										return v.toLocaleTimeString();
									}
								}],
								yAxes: [{
									display: true,
									ticks: {
										suggestedMin: 0
									}
								}],
							}
						}
					});
					this.graphDataLoaded = true;
				}
			);
		}
	}

}
