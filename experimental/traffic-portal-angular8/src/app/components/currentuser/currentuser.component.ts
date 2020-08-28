/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
import { Component, OnInit, OnDestroy } from '@angular/core';

import { Subscription } from 'rxjs';
import { first } from 'rxjs/operators';

import { Role, User } from '../../models';
import { APIService, AuthenticationService } from '../../services';

@Component({
	selector: 'currentuser',
	templateUrl: './currentuser.component.html',
	styleUrls: ['./currentuser.component.scss']
})
export class CurrentuserComponent implements OnInit, OnDestroy {

	currentUser: User;
	private subscription: Subscription;

	constructor(private readonly auth: AuthenticationService, private readonly api: APIService) {
	}

	ngOnInit() {
		this.currentUser = this.auth.currentUserValue;
		if (this.currentUser === null) {
			this.auth.updateCurrentUser().pipe(first()).subscribe(
				(r: boolean) => {
					if (r) {
						this.currentUser = this.auth.currentUserValue;
					}
				}
			);
		}
		this.subscription = this.auth.currentUser.subscribe(
			(u: User) => {
				this.currentUser = u;
			}
		);
	}

	ngOnDestroy() {
		this.subscription.unsubscribe();
	}
}
