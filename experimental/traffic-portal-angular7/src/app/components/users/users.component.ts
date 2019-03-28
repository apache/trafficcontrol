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

import { first } from 'rxjs/operators';

import { APIService, AuthenticationService } from '../../services';
import { User } from '../../models/user';

@Component({
	selector: 'users',
	templateUrl: './users.component.html',
	styleUrls: ['./users.component.scss']
})
export class UsersComponent implements OnInit {

	users: Array<User>;
	fuzzControl = new FormControl('');
	loading: boolean;
	myId: number;

	constructor (private readonly api: APIService, private readonly auth: AuthenticationService) {
		this.users = new Array<User>();
		this.loading = true;
		this.myId = -1;
	}

	ngOnInit () {
		// User may have navigated directly with a valid cookie - in which case current user is null
		if (this.auth.currentUserValue === null) {
			this.auth.updateCurrentUser().subscribe(
				v => {
					if (v) {
						this.myId = this.auth.currentUserValue.id;
					}
				}
			);
		} else {
			this.myId = this.auth.currentUserValue.id;
		}

		this.api.getUsers().pipe(first()).subscribe(
			r => {
				this.users = r;
				this.loading = false;
			}
		);
	}

	fuzzy (u: User): boolean {
		if (!this.fuzzControl.value) {
			return true;
		}
		const testVal = u.username.toLocaleLowerCase();
		let n = -1;
		for (const l of this.fuzzControl.value.toLocaleLowerCase()) {
			/* tslint:disable */
			if (!~(n = testVal.indexOf(l, n + 1))) {
			/* tslint:enable */
				return false;
			}
		}
		return true;
	}

}
