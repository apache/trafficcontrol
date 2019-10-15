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
import { first, map } from 'rxjs/operators';

import { Role, User } from '../models/user';
import { APIService } from './api.service';

@Injectable({ providedIn: 'root' })
export class AuthenticationService {
	private readonly currentUserSubject: BehaviorSubject<User>;
	public currentUser: Observable<User>;
	private readonly loggedInSubject: BehaviorSubject<boolean>;
	public loggedIn: Observable<boolean>;
	private readonly currentUserCapabilitiesSubject: BehaviorSubject<Array<string>>;
	public currentUserCapabilities: Observable<Array<string>>;

	constructor (private readonly http: HttpClient, private readonly api: APIService) {
		this.currentUserSubject = new BehaviorSubject<User>(null);
		this.loggedInSubject = new BehaviorSubject<boolean>(false);
		this.currentUserCapabilitiesSubject = new BehaviorSubject<Array<string>>([]);
		this.currentUser = this.currentUserSubject.asObservable();
		this.loggedIn = this.loggedInSubject.asObservable();
		this.currentUserCapabilities = this.currentUserCapabilitiesSubject.asObservable();
	}

	public get currentUserValue (): User {
		return this.currentUserSubject.value;
	}

	public get loggedInValue (): boolean {
		return this.loggedInSubject.value;
	}

	public get currentUserCapabilitiesValue (): Array<string> {
		return this.currentUserCapabilitiesSubject.value;
	}

	/**
	 * Updates the current user, and provides a way for callers to check if the update was succesful
	 * @returns An `Observable` which will emit a boolean value indicating the success of the update
	*/
	updateCurrentUser (): Observable<boolean> {
		return this.api.getCurrentUser().pipe(first()).pipe(map(
			(u: User) => {
				this.currentUserSubject.next(u);
				if (u.role) {
					this.api.getRoles(u.role).pipe(first()).pipe(map(
						(r: Role) => {
							this.currentUserCapabilitiesSubject.next(r.capabilities);
						}
					));
				}
				return true;
			},
			e => {
				console.error('Failed to update current user');
				console.debug('User update error: ', e);
				return false;
			}
		));
	}

	/**
	 * Logs in a user and, on succesful login, updates the current user.
	*/
	login (u: string, p: string): Observable<boolean> {
		return this.api.login(u, p).pipe(map(
			(resp) => {
				if (resp && resp.status === 200) {
					this.loggedInSubject.next(true);
					this.updateCurrentUser();
					return true;
				}
				return false;
			}
		));
	}

	logout () {
		this.currentUserSubject.next(null);
		this.loggedInSubject.next(false);
	}
}
