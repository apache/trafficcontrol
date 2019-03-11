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
import { HttpClient, HttpResponse } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { User } from '../models/user';
import { APIService } from './api.service';

@Injectable({ providedIn: 'root' })
export class AuthenticationService {
	private currentUserSubject: BehaviorSubject<User>;
	public currentUser: Observable<User>;
	private loggedInSubject: BehaviorSubject<boolean>;
	public loggedIn: Observable<boolean>;

	constructor(private http: HttpClient, private api: APIService) {
		this.currentUserSubject = new BehaviorSubject<User>(null);
		this.loggedInSubject = new BehaviorSubject<boolean>(false);
		this.currentUser = this.currentUserSubject.asObservable();
		this.loggedIn = this.loggedInSubject.asObservable();
	}

	public get currentUserValue(): User {
		return this.currentUserSubject.value;
	}

	public get loggedInValue(): boolean {
		return this.loggedInSubject.value;
	}

	private updateCurrentUser(): void {
		this.api.getCurrentUser().subscribe(
			r => {
				if (r.status === 200) {
					console.debug(r.body.response as User);
					this.currentUserSubject.next(r.body.response as User);
				}
			},
			e => {
				console.error("Failed to update current user");
			}
		);
	}

	login(u: string, p: string): Observable<boolean> {
		return this.api.login(u, p).pipe(map(
			(resp) => {
				if (resp && resp.status === 200) {
					this.loggedInSubject.next(true);
					this.updateCurrentUser();
					console.log("returning true");
					return true;
				}
				console.log("returning false");
				return false;
			}
		));
	}

	logout() {
		this.currentUserSubject.next(null);
	}
}
