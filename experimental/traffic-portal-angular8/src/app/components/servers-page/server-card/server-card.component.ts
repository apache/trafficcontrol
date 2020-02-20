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
import { Component, OnInit, Input } from '@angular/core';

import { Server } from '../../../models';

@Component({
	selector: 'server-card',
	templateUrl: './server-card.component.html',
	styleUrls: ['./server-card.component.scss']
})
export class ServerCardComponent implements OnInit {

	@Input() server: Server;

	constructor() { }

	ngOnInit(): void {
	}

	public down(): boolean {
		return this.server.status === 'ADMIN_DOWN' || this.server.status === 'OFFLINE';
	}

	/**
	 * cacheServer returns 'true' if this component's server is a 'cache server', 'false' otherwise.
	*/
	public cacheServer(): boolean {
		return this.server.type && (this.server.type.indexOf('EDGE') === 0 || this.server.type.indexOf('MID') === 0);
	}

}
