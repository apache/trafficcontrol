import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpResponse } from '@angular/common/http';
import { BehaviorSubject, Observable, throwError } from 'rxjs';
import { map, first, catchError } from 'rxjs/operators';

import { User } from '../models/user';

@Injectable({ providedIn: 'root' })
export class APIService {
	public API_VERSION = '1.5';

	private cookies: string;

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
		const options = {headers: new HttpHeaders({'Content-Type': 'application/json',
		                                           'Cookie': this.cookies ? this.cookies : ''}),
		                 observe: 'response' as 'response',
		                 responseType: 'json' as 'json',
		                 body: data};
		const ret = this.http.request(method, path, options)/*.pipe(map((response) => {
			//TODO pass alerts to the alert service
			// (TODO create the alert service)
			console.log("http.request made");
			let resp = response as HttpResponse<any>;
			if ('Cookie' in resp.headers) {
				this.cookies = resp.headers['Cookie']
			}
			console.log("returning resp");
			return resp;
		})).pipe(catchError((e, caught) => {
			console.log("wtf?: ", e);
			return throwError(e);
		}));*/
		ret.subscribe(r => {
			console.log("got a response: ", r);
		},
		e => {
			console.error("got an error: ", e)
		});
		return ret;
	}

	public login(u, p): Observable<HttpResponse<any>> {
		const path = 'http://localhost:4000/api/'+this.API_VERSION+'/user/login';
		return this.post(path, {u, p});
	}

	public getCurrentUser(): Observable<HttpResponse<any>> {
		const path = '/api/'+this.API_VERSION+'/user/current';
		return this.get(path);
	}
}
