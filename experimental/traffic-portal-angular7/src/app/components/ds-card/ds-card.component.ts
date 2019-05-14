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
	midBandwidth: Array<number>;
	edgeBandwidth: Array<number>;
	labels: Array<Date>;

	// Need this to access merged namespace for string conversions
	Protocol = Protocol;

	chartOptions: any;

	private loaded: boolean;
	public graphDataLoaded: boolean;

	constructor (private readonly api: APIService, private readonly elementRef: ElementRef) {
		this.available = 100;
		this.maintenance = 0;
		this.utilized = 0;
		this.loaded = false;
		this.edgeBandwidth = new Array<number>();
		this.midBandwidth = new Array<number>();
		this.labels = new Array<Date>();
		this.graphDataLoaded = false;
		this.chartOptions = {
			type: 'line',
			data: {
				labels: [],
				datasets: []
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
						callback: (v, unused_i, unused_values) => {
							return v.toLocaleTimeString();
						}
					}],
					yAxes: [{
						display: true,
						ticks: {
							suggestedMin: 0
						}
					}]
				}
			}
		};
	}

	/**
	 * Handles the destruction of a chart and all of its constituent data. Does nothing if
	 * `this.chart` is `null` or `undefined`.
	*/
	private destroyChart () {
		if (this.chart) {
			this.chart.destroy();
			this.chart = null;
			this.graphDataLoaded = false;
			this.chartOptions.data = {datasets: [], labels: []};
		}
	}

	/**
	 * Triggered when the details element is opened or closed, and fetches the latest stats.
	 * @param e A DOM Event caused then the open/close state of a `<details>` element changes.
	 *
	 * this will only fetch health and capacity information once per page load, but will update the
	 * Bandwidth graph every time the details panel is opened. Bandwidth data is calculated using
	 * 60s intervals starting at 00:00 UTC the current day and ending at the current date/time.
	*/
	toggle (e: Event) {
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
							if (r.totalOnline === 0) {
								this.healthy = 0;
							} else {
								this.healthy = Number(r.totalOnline) / (Number(r.totalOnline) + Number(r.totalOffline));
							}
						}
					}
				);
			}
			const now = new Date();
			now.setUTCMilliseconds(0); // apparently `const` doesn't care about this
			const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
			const canvas = document.getElementById('canvas-' + String(this.deliveryService.id)) as HTMLCanvasElement;
			this.chart = new Chart(canvas, this.chartOptions);

			this.api.getDSKBPS(this.deliveryService.xmlId, today, now, '60s').pipe(first()).subscribe(
				data => {
					if (this.chart === null) {
						return;
					}
					if (data === undefined || data === null || data.series === undefined || data.series === null || data.series.values === undefined || data.series.values === null) {
						this.destroyChart();

						const ctx = canvas.getContext('2d');
						ctx.font = '30px serif';
						ctx.fillStyle = 'black';
						ctx.textAlign = 'center';
						ctx.fillText('No Data', canvas.width / 2., canvas.height / 2.);
						this.graphDataLoaded = true;
						return;
					}
					for (const d of data.series.values) {
						this.chart.data.labels.push(new Date(d[0]));
						this.edgeBandwidth.push(d[1]);
					}
					this.chart.data.datasets.push({
						label: 'Edge-tier Bandwidth',
						data: this.edgeBandwidth,
						borderColor: '#3CBA9F',
						fill: true,
						fillColor: '#3CBA9F'
					});
					this.graphDataLoaded = true;
					this.chart.update();
				}
			);

			this.api.getDSKBPS(this.deliveryService.xmlId, today, now, '60s', true).pipe(first()).subscribe(
				data => {
					if (this.chart === null) {
						return;
					}
					if (data === undefined || data === null || data.series === undefined || data.series === null || data.series.values === undefined || data.series.values === null) {
						this.destroyChart();

						const ctx = canvas.getContext('2d');
						ctx.font = '30px serif';
						ctx.fillStyle = 'black';
						ctx.textAlign = 'center';
						ctx.fillText('No Data', canvas.width / 2., canvas.height / 2.);
						this.graphDataLoaded = true;
						return;
					}
					for (const d of data.series.values) {
						this.midBandwidth.push(d[1]);
					}
					this.chart.data.datasets.push({
						label: 'Mid-tier Bandwidth',
						data: this.midBandwidth,
						borderColor: '#BA3C57',
						fill: true,
						fillColor: '#BA3C57'
					});
					this.graphDataLoaded = true;
					this.chart.update();
				}
			);
		} else if (this.chart) {
			this.destroyChart();
		}
	}

}
