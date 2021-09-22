import { Injectable } from "@angular/core";
import {
	ActivatedRouteSnapshot,
	CanActivate,
	RouterStateSnapshot,
	UrlTree
} from "@angular/router";
import { Observable } from "rxjs";
import {AuthenticationService} from "./shared/authentication/authentication.service";

/**
 *
 */
@Injectable()
export class AuthenticationGuard implements CanActivate {
	constructor(private readonly auth: AuthenticationService) {
	}

	/**
	 *
	 * @param route
	 * @param state
	 * @returns boolean
	 */
	public canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
		console.log(this.auth.currentUser);
		return true;
	}
}
