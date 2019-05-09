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

import { BehaviorSubject, Observable } from 'rxjs';

import { APIService } from '../../services';
import { DeliveryService } from '../../models/deliveryservice';
import { DataPoint } from '../../models/datapoint';

@Component({
	selector: 'deliveryservice',
	templateUrl: './deliveryservice.component.html',
	styleUrls: ['./deliveryservice.component.scss']
})
export class DeliveryserviceComponent implements OnInit {

	deliveryservice = new DeliveryService();
	loaded = new Map([['main', false], ['bandwidth', false]]);

	private readonly bandwidthDataSubject: BehaviorSubject<Array<Array<DataPoint>>>;
	public bandwidthData: Observable<Array<Array<DataPoint>>>;

	midBandwidth: Array<number>;
	edgeBandwidth: Array<DataPoint>;
	labels: Array<Date>;

	to: Date;
	from: Date;
	fromDate: FormControl;
	fromTime: FormControl;
	toDate: FormControl;
	toTime: FormControl;

	constructor(private readonly route: ActivatedRoute, private readonly api: APIService, private readonly element: ElementRef<HTMLCanvasElement>) {
		this.labels = new Array<Date>();
		this.midBandwidth = new Array<number>();
		this.edgeBandwidth = new Array<DataPoint>();
		this.bandwidthData = this.bandwidthDataSubject.asObservable();
	}

	public get bandwidthDataValue (): Array<Array<DataPoint>> {
		return this.bandwidthDataSubject.value;
	}

	ngOnInit() {
		const DSID = this.route.snapshot.paramMap.get('id');

		this.to = new Date();
		this.to.setUTCMilliseconds(0);
		this.from = new Date(this.to.getFullYear(), this.to.getMonth(), this.to.getDate());
		const dateStr = String(this.from.getFullYear()).padStart(4, '0').concat('-', String(this.from.getMonth()).padStart(2, '0').concat('-', String(this.from.getDate()).padStart(2, '0')));
		this.fromDate = new FormControl(dateStr);
		this.fromTime = new FormControl("00:00");
		this.toDate = new FormControl(dateStr);
		const timeStr = String(this.to.getHours()).padStart(2, '0').concat(':', String(this.to.getMinutes()).padStart(2, '0'))
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
		this.api.getDSKBPS(this.deliveryservice.xmlId, this.from, this.to, '300s').subscribe(
			va => {
				if (va === null || va === undefined){
					this.edgeBandwidth = null;
					return;
				}

				for (const v of va) {
					this.edgeBandwidth.push({t: new Date(v[0]), y: v[1]});
				}
				this.bandwidthDataSubject.next([this.edgeBandwidth]);
				console.log("loaded graph", this.edgeBandwidth);
			}
		);
	}

}
