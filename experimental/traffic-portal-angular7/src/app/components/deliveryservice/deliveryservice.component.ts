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
import { DeliveryService } from '../../models/deliveryservice';

@Component({
	selector: 'deliveryservice',
	templateUrl: './deliveryservice.component.html',
	styleUrls: ['./deliveryservice.component.scss']
})
export class DeliveryserviceComponent implements OnInit {

	deliveryservice = new DeliveryService();
	loading = false;

	constructor(private route: ActivatedRoute, private api: APIService) { }

	ngOnInit() {
		const DSID = this.route.snapshot.paramMap.get("id");
		this.api.getDeliveryServices(parseInt(DSID)).subscribe(
			(d: DeliveryService) => {
				this.deliveryservice = d;
				this.loading = true;
			}
		);
	}

}
