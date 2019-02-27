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

	constructor(private http: HttpClient, private api: APIService) {
		this.currentUserSubject = new BehaviorSubject<User>(null);
		this.currentUser = this.currentUserSubject.asObservable();
	}

	public get currentUserValue(): User {
		return this.currentUserSubject.value;
	}

	private updateCurrentUser(): void {
		// this.api.getCurrentUser().subscribe(
		// 	r => {
		// 		if (r.status === 200) {
		// 			console.debug(r.body.response as User);
		// 			this.currentUserSubject.next(r.body.response as User);
		// 		}
		// 	},
		// 	e => {
		// 		console.error("Failed to update current user");
		// 	}
		// );
	}

	login(u: string, p: string): Observable<boolean> {
		return this.api.login(u, p).pipe(map(
			(resp) => {
				if (resp && resp.status === 200) {
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
