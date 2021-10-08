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
import { Injectable } from "@angular/core";
import {
	ActivatedRouteSnapshot,
	CanActivate, CanLoad, Route,
	RouterStateSnapshot, UrlSegment,
	UrlTree
} from "@angular/router";
import { Observable } from "rxjs";
import {AuthenticationService} from "../shared/authentication/authentication.service";

/**
 * AuthenticationGuard ensures that the user is logged in.
 */
@Injectable()
export class AuthenticatedGuard implements CanActivate, CanLoad {
	constructor(private readonly auth: AuthenticationService) {
	}

	/**
	 * canActivate determines whether or not a user can activate an already loaded route.
	 *
	 * @param route Route snapshot.
	 * @param state Route state snapshot.
	 * @returns many
	 */
	public canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean | UrlTree>
	| Promise<boolean | UrlTree> | boolean | UrlTree {
		return this.auth.currentUser !== null;
	}

	/**
	 * canLoad determines whether or not the current user can load/request the given route.
	 *
	 * @param route Requested route.
	 * @param segments URL segments.
	 * @returns many
	 */
	public canLoad(route: Route, segments: UrlSegment[]): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
		return this.auth.currentUser !== null;
	}
}
