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

import { Server } from '../../models';

@Injectable({providedIn: 'root'})
export class ServerService extends APIService {
	public getServers(): Observable<Array<Server>> {
		let path = `/api/${this.API_VERSION}/servers`;
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<Server>;
			}
		));
	}

	constructor(http: HttpClient) {
		super(http);
	}
}
