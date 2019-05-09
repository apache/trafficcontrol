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
import { ActivatedRoute } from '@angular/router';

import { Chart } from 'chart.js';

import { APIService } from '../../services';
import { DeliveryService } from '../../models/deliveryservice';

@Component({
	selector: 'deliveryservice',
	templateUrl: './deliveryservice.component.html',
	styleUrls: ['./deliveryservice.component.scss']
})
export class DeliveryserviceComponent implements OnInit {

	deliveryservice = new DeliveryService();
	loaded = new Map([['main', false], ['bandwidth', false]]);

	chart: Chart;
	chartOptions: any;
	midBandwidth: Array<number>;
	edgeBandwidth: Array<number>;
	labels: Array<Date>;

	now: Date;
	today: Date;
	fromDate: FormControl;
	fromTime: FormControl;
	toDate: FormControl;
	toTime: FormControl;

	constructor(private readonly route: ActivatedRoute, private readonly api: APIService) {
		this.labels = new Array<Date>();
		this.midBandwidth = new Array<number>();
		this.edgeBandwidth = new Array<number>();
		this.chartOptions = {
			type: 'line',
			data: {
				labels: [],
				datasets: []
			},
			options: {
				legend: {
					display: true
				},
				title: {
					display: true,
					text: "Bandwidth of Cache Tiers"
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

	ngOnInit() {
		const DSID = this.route.snapshot.paramMap.get('id');

		this.now = new Date();
		this.now.setUTCMilliseconds(0);
		this.today = new Date(this.now.getFullYear(), this.now.getMonth(), this.now.getDate());
		const dateStr = String(this.today.getFullYear()).padStart(4, '0').concat('-', String(this.today.getMonth()).padStart(2, '0').concat('-', String(this.today.getDate()).padStart(2, '0')));
		this.fromDate = new FormControl(dateStr);
		this.fromTime = new FormControl("00:00");
		this.toDate = new FormControl(dateStr);
		const timeStr = String(this.now.getHours()).padStart(2, '0').concat(':', String(this.now.getMinutes()).padStart(2, '0'))
		console.log("Loading timeStr: ", timeStr);
		this.toTime = new FormControl(timeStr);


		this.api.getDeliveryServices(parseInt(DSID)).subscribe(
			(d: DeliveryService) => {
				this.deliveryservice = d;
				this.loaded['main'] = true;
				this.loadGraph();
			}
		);

	}

	loadGraph() {
		const now = new Date();
		now.setUTCMilliseconds(0);
		console.log("loaded graph");
	}

}
