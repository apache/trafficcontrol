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

import { Subject } from 'rxjs';
import { first } from 'rxjs/operators';

import { APIService } from '../../services';
import { DeliveryService, Protocol } from '../../models/deliveryservice';
import { DataPoint, DataSet } from '../../models/data';

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
	chartData: Subject<Array<DataSet>>;
	midBandwidthData: DataSet;
	edgeBandwidthData: DataSet;

	// Need this to access merged namespace for string conversions
	Protocol = Protocol;

	private loaded: boolean;
	public graphDataLoaded: boolean;

	constructor (private readonly api: APIService) {
		this.available = 100;
		this.maintenance = 0;
		this.utilized = 0;
		this.loaded = false;
		this.chartData = new Subject<Array<DataSet>>();

		this.edgeBandwidthData = {
			label: 'Edge-tier Bandwidth',
			data: new Array<DataPoint>(),
			borderColor: '#3CBA9F',
			fill: true,
			fillColor: '#3CBA9F'
		} as DataSet;

		this.midBandwidthData = {
			label: 'Mid-tier Bandwidth',
			data: new Array<DataPoint>(),
			borderColor: '#BA3C57',
			fill: true,
			fillColor: '#BA3C57'
		} as DataSet;

		this.graphDataLoaded = false;
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
			now.setUTCMilliseconds(0);
			const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
			this.loadChart(now, today);
		}
	}

	private loadChart(now: Date, today: Date) {
		this.api.getDSKBPS(this.deliveryService.xmlId, today, now, '60s').pipe(first()).subscribe(
			data => {
				if (data === null || data.series === undefined || data.series === null || data.series.values === undefined || data.series.values === null) {
					this.graphDataLoaded = true;
					this.chartData.next([null, this.midBandwidthData])
					return;
				}
				for (const d of data.series.values) {
					if (d[1] === null) {
						continue;
					}
					this.edgeBandwidthData.data.push({t: new Date(d[0]), y: d[1]} as DataPoint);
				}
				this.chartData.next([this.edgeBandwidthData, this.midBandwidthData]);
				this.graphDataLoaded = true;
			}
		);

		this.api.getDSKBPS(this.deliveryService.xmlId, today, now, '60s', true).pipe(first()).subscribe(
			data => {
				if (data === null || data.series === undefined || data.series === null || data.series.values === undefined || data.series.values === null) {
					this.chartData.next([this.edgeBandwidthData, null])
					this.graphDataLoaded = true;
					return;
				}
				for (const d of data.series.values) {
					if (d[1] === null) {
						continue;
					}
					this.midBandwidthData.data.push({t: new Date(d[0]), y: d[1]} as DataPoint);
				}
				this.chartData.next([this.edgeBandwidthData, this.midBandwidthData]);
				this.graphDataLoaded = true;
			}
		);
	}
}
