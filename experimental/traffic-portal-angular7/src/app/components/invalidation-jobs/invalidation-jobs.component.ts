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
import { ActivatedRoute } from '@angular/router';

import { APIService } from '../../services';
import { DeliveryService } from '../../deliveryservice';
import { InvalidationJob } from '../../models/invalidation';

@Component({
	selector: 'invalidation-jobs',
	templateUrl: './invalidation-jobs.component.html',
	styleUrls: ['./invalidation-jobs.component.scss']
})
export class InvalidationJobsComponent implements OnInit {

	deliveryservice: DeliveryService = {};
	jobs: Array<InvalidationJob>;
	xmlid: string;

	constructor (private readonly route: ActivatedRoute, private readonly api: APIService) { }

	ngOnInit () {
		const id = parseInt(this.route.snapshot.paramMap.get('id'), 10);
		this.api.getInvalidationJobs({dsId: id}).subscribe(
			r => {
				this.jobs = r;
				console.log(r);
				console.log(this.jobs);
			}
		);
		this.api.getDeliveryServices(id).subscribe(
			r => {
				this.deliveryservice = r;
			}
		);
	}

}
