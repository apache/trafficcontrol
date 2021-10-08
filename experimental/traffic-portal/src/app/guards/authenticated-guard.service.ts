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
 *
 */
@Injectable()
export class AuthenticatedGuard implements CanActivate, CanLoad {
	constructor(private readonly auth: AuthenticationService) {
	}

	public canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
		return this.auth.currentUser !== null;
	}

	public canLoad(route: Route, segments: UrlSegment[]): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
		return this.auth.currentUser !== null;
	}
}
