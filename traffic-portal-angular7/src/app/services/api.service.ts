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
import { HttpClient, HttpHeaders, HttpResponse } from '@angular/common/http';
import { BehaviorSubject, Observable, throwError } from 'rxjs';
import { map, first, catchError } from 'rxjs/operators';

import { DeliveryService } from '../models/deliveryservice';
import { User } from '../models/user';

@Injectable({ providedIn: 'root' })
/**
 * The APIService provides access to the Traffic Ops API. Its methods should be kept API-version
 * agnostic (from the caller's perspective), and always return `Observable`s.
*/
export class APIService {
	public API_VERSION = '1.5';

	// private cookies: string;

	constructor(private http: HttpClient) {

	}

	private delete(path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('delete', path, data);
	}
	private get(path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('get', path, data);
	}
	private head(path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('head', path, data);
	}
	private options(path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('options', path, data);
	}
	private patch(path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('patch', path, data);
	}
	private post(path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('post', path, data);
	}
	private push(path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('push', path, data);
	}

	private do(method: string, path: string, data?: Object): Observable<HttpResponse<any>> {
		const options = {headers: new HttpHeaders({'Content-Type': 'application/json'}),
		                 observe: 'response' as 'response',
		                 responseType: 'json' as 'json',
		                 body: data};
		return this.http.request(method, path, options).pipe(map((response) => {
			//TODO pass alerts to the alert service
			// (TODO create the alert service)
			return response as HttpResponse<any>;
		}));
	}

	/**
	 * Performs authentication with the Traffic Ops server.
	 * @param u The username to be used for authentication
	 * @param p The password of user `u`
	 * @returns An observable that will emit the entire HTTP response
	*/
	public login(u:string, p:string): Observable<HttpResponse<any>> {
		const path = 'http://localhost:4000/api/'+this.API_VERSION+'/user/login';
		return this.post(path, {u, p});
	}

	/**
	 * Fetches the current user from Traffic Ops
	 * @returns An observable that will emit a `User` object representing the current user.
	*/
	public getCurrentUser(): Observable<User> {
		const path = '/api/'+this.API_VERSION+'/user/current';
		return this.get(path).pipe(map(
			r => {
				return r.body.response as User;
			}
		));
	}

	/**
	 * Gets a list of all visible Delivery Services
	 * @returns An observable that will emit an array of `DeliveryService` objects.
	*/
	public getDeliveryServices(): Observable<DeliveryService[]> {
		const path = '/api/'+this.API_VERSION+'/deliveryservices';
		return this.get(path).pipe(map(
			r => {
				return r.body.response as DeliveryService[];
			}
		));
	}
}
