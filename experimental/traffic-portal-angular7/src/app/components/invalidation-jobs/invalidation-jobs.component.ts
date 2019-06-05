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

import { Subject } from 'rxjs';

import { APIService } from '../../services';
import { DeliveryService } from '../../models/deliveryservice';
import { InvalidationJob } from '../../models/invalidation';

@Component({
	selector: 'invalidation-jobs',
	templateUrl: './invalidation-jobs.component.html',
	styleUrls: ['./invalidation-jobs.component.scss']
})
export class InvalidationJobsComponent implements OnInit {

	deliveryservice: DeliveryService ;
	jobs: Array<InvalidationJob>;
	now: Date;
	showDialog: Subject<boolean>;

	regexp = new FormControl('/');
	ttl = new FormControl(178);
	startDate = new FormControl('');
	startTime = new FormControl('');


	constructor (private readonly route: ActivatedRoute, private readonly api: APIService) {
		this.deliveryservice = {active: true} as DeliveryService;
		this.jobs = new Array<InvalidationJob>();
		this.showDialog = new Subject<boolean>();
	}

	ngOnInit () {
		this.now = new Date();
		const id = parseInt(this.route.snapshot.paramMap.get('id'), 10);
		this.api.getInvalidationJobs({dsId: id}).subscribe(
			r => {
				// The values returned by the API are not RFC-compliant at the time of this writing,
				// so we need to do some pre-processing on them.
				for (const j of r) {
					const tmp = Array.from(String(j.startTime).split(' ').join('T'));
					tmp.splice(-3, 3);
					j.startTime = new Date(tmp.join(''));
					this.jobs.push(j);
				}
			}
		);
		this.api.getDeliveryServices(id).subscribe(
			(r: DeliveryService) => {
				this.deliveryservice = r;
			}
		);
	}

	public endDate (j: InvalidationJob): Date {
		const tmp = j.parameters.split(':');
		if (tmp.length !== 2) {
			throw new Error('Malformed job parameters: "' + j.parameters + '" (id: ' + String(j.id) + ')');
		}
		const ttl = parseInt(tmp[1], 10);
		if (isNaN(ttl)) {
			throw new Error('Invalid TTL: "' + tmp[1] + '" (job id: ' + String(j.id) + ')');
		}
		return new Date(new Date(j.startTime.getTime() + ttl*60*60*1000));
	}

	public newJob (e?: Event) {
		if (e) {
			e.preventDefault();
		}
		this.showDialog.next(true);
	}

	public closeDialog (e?: Event) {
		if (e) {
			e.preventDefault();
		}


		let job: InvalidationJob;
		try {
			job = {
				dsId: this.deliveryservice.id,
				parameters: 'TTL:' + String(this.ttl.value),
				regex: (new RegExp(this.regexp.value)).toString().replace(/^\/|\/$/g, '').replace('\\/', '/'),
				startTime: this.startDate.value.concat(' ', this.startTime.value + ':00'),
				ttl: this.ttl.value
			};
		} catch (err) {
			console.error(err);
			return;
		}

		this.api.createInvalidationJob(job).subscribe(
			r => {
				if (r) {
					this.api.getInvalidationJobs({dsId: this.deliveryservice.id}).subscribe(
						r => {
							this.jobs = new Array<InvalidationJob>();
							for (const j of r) {
								const tmp = Array.from(String(j.startTime).replace(' ', 'T'));
								tmp.splice(-3, 3);
								j.startTime = new Date(tmp.join(''));
								this.jobs.push(j);
							}
						}
					);
				} else {
					console.warn("failure");
				}
			},
			e => {
				console.error("error: ", e);
			}
		)
	}

}
