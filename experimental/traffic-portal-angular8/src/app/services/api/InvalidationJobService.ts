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

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { APIService } from './apiservice';

import { DeliveryService, InvalidationJob, User } from '../../models';

@Injectable({providedIn: 'root'})
export class InvalidationJobService extends APIService {
	public getInvalidationJobs (opts?: {id: number} |
	                                   {userId: number} |
	                                   {user: User} |
	                                   {dsId: number} |
	                                   {deliveryService: DeliveryService}): Observable<Array<InvalidationJob>> {
		let path = '/api/' + this.API_VERSION + '/jobs';
		if (opts) {
			path += '?';
			if (opts.hasOwnProperty('id')) {
				path += 'id=' + String((opts as {id: number}).id);
			} else if (opts.hasOwnProperty('dsId')) {
				path += 'dsId=' + String((opts as {dsId: number}).dsId);
			} else if (opts.hasOwnProperty('userId')) {
				path += 'userId=' + String((opts as {userId: number}).userId);
			} else if (opts.hasOwnProperty('deliveryService')) {
				path += 'dsId=' + String((opts as {deliveryService: DeliveryService}).deliveryService.id);
			} else {
				path += 'userId=' + String((opts as {user: User}).user.id);
			}
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<InvalidationJob>;
			}
		));
	}

	public createInvalidationJob (job: InvalidationJob): Observable<boolean> {
		const path = '/api/' + this.API_VERSION + '/user/current/jobs';
		return this.post(path, job).pipe(map(
			r => true,
			e => false
		));
	}

	constructor(http: HttpClient) {
		super(http);
	}
}
