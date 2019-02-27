import { Injectable } from '@angular/core';
import { HttpRequest, HttpHandler, HttpEvent, HttpInterceptor } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

import { AuthenticationService } from '../services';

@Injectable()
export class ErrorInterceptor implements HttpInterceptor {
	constructor(private authenticationService: AuthenticationService) {}

	intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
		return next.handle(request).pipe(catchError(err => {
			// if (err.status === 401 || err.status === 403) {
			// 	this.authenticationService.logout();
			// 	location.reload(true);
			// }

			// console.error("Error response: ", err);
			// const error = err.error && err.error.message ? err.error.message : err.statusText;
			return throwError(err);
		}));
	}
}
