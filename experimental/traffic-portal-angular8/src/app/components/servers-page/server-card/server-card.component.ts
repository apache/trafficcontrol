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
import { first } from 'rxjs/operators';

import { faTimesCircle, faCheckCircle, faClock } from '@fortawesome/free-solid-svg-icons';

import { Server, Servercheck, checkMap } from '../../../models';
import { ServerService } from '../../../services/api';

@Component({
	selector: 'server-card',
	templateUrl: './server-card.component.html',
	styleUrls: ['./server-card.component.scss']
})
export class ServerCardComponent implements OnInit {

	@Input() server: Server;

	checks: Map<string, number|boolean>;
	open: boolean;

	constructor(private readonly api: ServerService) {
		this.open = false;
	}

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

	public updPendingIcon() {
		return this.server.updPending ? faClock : faCheckCircle;
	}

	public updPendingTitle() {
		return this.server.updPending ? "Updates are pending" : "No updates are pending";
	}

	public revalPendingIcon() {
		return this.server.revalPending ? faClock : faCheckCircle;
	}

	public revalPendingTitle() {
		return this.server.revalPending ? "Revalidations are pending" : "No revalidations are pending";
	}

	public toggle(e: Event): void {
		this.open = !this.open;
		if (this.open) {
			this.api.getServerChecks(this.server.id).pipe(first()).subscribe(
				(s: Servercheck) => {
					this.checks = s.checkMap()
				}
			);
		}
	}
}
