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
import { Component, ElementRef, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';

import { Subject } from 'rxjs';

import { AlertService, APIService } from '../../services';
import { DeliveryService } from '../../models/deliveryservice';
import { DataPoint, DataSet, TPSData } from '../../models/data';

@Component({
	selector: 'deliveryservice',
	templateUrl: './deliveryservice.component.html',
	styleUrls: ['./deliveryservice.component.scss']
})
export class DeliveryserviceComponent implements OnInit {

	deliveryservice = new DeliveryService();
	loaded = new Map([['main', false], ['bandwidth', false]]);

	bandwidthData: Subject<Array<DataSet>>;
	TPSChartData: Subject<Array<DataSet>>;

	edgeBandwidth: DataSet;
	midBandwidth: DataSet;

	edgeTPSData: DataSet;
	midTPSData: DataSet;

	to: Date;
	from: Date;
	fromDate: FormControl;
	fromTime: FormControl;
	toDate: FormControl;
	toTime: FormControl;

	bucketSize = 300; // seconds

	constructor(private readonly route: ActivatedRoute, private readonly api: APIService, private readonly alerts: AlertService) {
		this.midBandwidth = {label: "Mid-Tier", borderColor: "#3CBA9F", fill: false, backgroundColor: "#3CBA9F", data: new Array<DataPoint>()} as DataSet;
		this.edgeBandwidth = {label: "Edge-Tier", borderColor: "#BA3C57", fill: false, backgroundColor: "#BA3C57", data: new Array<DataPoint>()} as DataSet;
		this.edgeTPSData = {label: "Edge-Tier", borderColor: "#BA3C57", fill: false, backgroundColor: "#BA3C57", data: new Array<DataPoint>()} as DataSet;
		this.midTPSData = {label: "Mid-Tier", borderColor: "#3CBA9F", fill: false, backgroundColor: "#3CBA9F", data: new Array<DataPoint>()} as DataSet;
		this.bandwidthData = new Subject<Array<DataSet>>();
		this.TPSChartData = new Subject<Array<DataSet>>();
	}

	ngOnInit() {
		const DSID = this.route.snapshot.paramMap.get('id');

		this.to = new Date();
		this.to.setUTCMilliseconds(0);
		this.from = new Date(this.to.getFullYear(), this.to.getMonth(), this.to.getDate());
		const dateStr = String(this.from.getFullYear()).padStart(4, '0').concat('-', String(this.from.getMonth() + 1).padStart(2, '0').concat('-', String(this.from.getDate()).padStart(2, '0')));
		this.fromDate = new FormControl(dateStr);
		this.fromTime = new FormControl("00:00");
		this.toDate = new FormControl(dateStr);
		const timeStr = String(this.to.getHours()).padStart(2, '0').concat(':', String(this.to.getMinutes()).padStart(2, '0'))
		this.toTime = new FormControl(timeStr);

		this.api.getDeliveryServices(parseInt(DSID)).subscribe(
			(d: DeliveryService) => {
				this.deliveryservice = d;
				this.loaded['main'] = true;
				this.loadBandwidth();
				this.loadTPS();
			}
		);
	}

	newDateRange() {
		this.to = new Date(this.toDate.value.concat('T', this.toTime.value));
		this.from = new Date(this.fromDate.value.concat('T', this.fromTime.value));

		// I need to set these explicitly, just in case - the API can't handle millisecond precision
		this.to.setUTCMilliseconds(0);
		this.from.setUTCMilliseconds(0);
		this.bucketSize = 60;
		this.loadBandwidth();
		this.loadTPS();
	}

	loadBandwidth() {
		// Edge-tier data
		this.api.getDSKBPS(this.deliveryservice.xmlId, this.from, this.to, String(this.bucketSize) + 's', false).subscribe(
			data => {
				const va = new Array<DataPoint>();
				for (const v of data.series.values) {
					if (v[1] === null) {
						continue;
					}
					va.push({t: new Date(v[0]), y: v[1]} as DataPoint);
				}
				this.edgeBandwidth.data = va;
				this.bandwidthData.next([this.edgeBandwidth, this.midBandwidth]);
			},
			(e: Error) => {
				this.alerts.newAlert("warning", "Edge-Tier bandwidth data not found!");
				console.debug(e);
			}
		);

		// Mid-tier data
		this.api.getDSKBPS(this.deliveryservice.xmlId, this.from, this.to, String(this.bucketSize) + 's', true).subscribe(
			data => {
				const va = new Array<DataPoint>();
				for (const v of data.series.values) {
					if (v[1] === null) {
						continue;
					}
					va.push({t: new Date(v[0]), y: v[1]} as DataPoint);
				}
				this.midBandwidth.data = va;
				this.bandwidthData.next([this.edgeBandwidth, this.midBandwidth]);
			},
			(e: Error) => {
				this.alerts.newAlert("warning", "Mid-Tier bandwidth data not found!");
				console.debug(e);
			}
		);
	}

	loadTPS() {
		// Edge-tier data
		// this.api.getDSTPS(this.deliveryservice.xmlId, this.from, this.to, String(this.bucketSize) + 's').subscribe(
		// 	data => {
		// 		if (data === null || data.series === undefined || data.series === null || data.series.values === undefined || data.series.values === null) {
		// 			this.alerts.newAlert("warning", "Edge-Tier transaction data not found!");
		// 			return;
		// 		}

		// 		const va = new Array<DataPoint>();
		// 		for (const v of data.series.values) {
		// 			if (v[1] === null) {
		// 				continue;
		// 			}
		// 			va.push({t: new Date(v[0]), y: v[1]} as DataPoint);
		// 		}
		// 		this.edgeTPSData.data = va;
		// 		this.TPSChartData.next([this.edgeTPSData, this.midTPSData]);
		// 	}
		// );

		this.api.getAllDSTPSData(this.deliveryservice.xmlId, this.from, this.to, String(this.bucketSize) + 's', false).subscribe(
			(data: TPSData) => {
				data.total.dataSet.label = "Total (edge-tier)";
				data.total.dataSet.borderColor = "#3C96BA";
				data.total.dataSet.borderDash = [5,15];
				data.success.dataSet.label = "Successful Responses (edge-tier)";
				data.success.dataSet.borderColor = "#3CBA5F";
				data.success.dataSet.borderDash = [5, 15];
				data.redirection.dataSet.label = "Redirection Responses (edge-tier)";
				data.redirection.dataSet.borderColor = "#9f3CBA";
				data.redirection.dataSet.borderDash = [5, 15];
				data.clientError.dataSet.label = "Client Error Responses (edge-tier)";
				data.clientError.dataSet.borderColor = "#BA9E3B";
				data.clientError.dataSet.borderDash = [5, 15];
				data.serverError.dataSet.label = "Server Error Responses (edge-tier)";
				data.serverError.dataSet.borderColor = "#BA3C57";
				data.serverError.dataSet.borderDash = [5, 15];

				this.TPSChartData.next([
					data.total.dataSet,
					data.success.dataSet,
					data.redirection.dataSet,
					data.clientError.dataSet,
					data.serverError.dataSet
				]);
			},
			(e: Error) => {
				console.debug(e);
				this.alerts.newAlert("warning", "Edge-Tier transaction data not found!");
			}
		);

		// this.api.getDSTPS(this.deliveryservice.xmlId, this.from, this.to, String(this.bucketSize) + 's').subscribe(
		// 	data => {
		// 		if (data === null || data.series === undefined || data.series === null || data.series.values === undefined || data.series.values === null) {
		// 			this.alerts.newAlert("warning", "Mid-Tier transaction data not found!");
		// 			return;
		// 		}

		// 		const va = new Array<DataPoint>();
		// 		for (const v of data.series.values) {
		// 			if (v[1] === null) {
		// 				continue;
		// 			}
		// 			va.push({t: new Date(v[0]), y: v[1]} as DataPoint);
		// 		}
		// 		this.midTPSData.data = va;
		// 		this.TPSChartData.next([this.edgeTPSData, this.midTPSData]);
		// 	}
		// );
	}

}
