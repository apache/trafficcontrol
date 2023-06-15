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

import { CurrentUserService } from "src/app/shared/current-user/current-user.service";

/**
 * AuthenticationGuard ensures that the user is logged in.
 */
@Injectable()
export class AuthenticatedGuard  {
	constructor(private readonly auth: CurrentUserService) {
	}

	/**
	 * canActivate determines whether or not a user can activate an already loaded route.
	 *
	 * @returns Whether or not the route can be activated.
	 */
	public async canActivate(): Promise<boolean>  {
		return this.auth.fetchCurrentUser();
	}

	/**
	 * canLoad determines whether or not the current user can load/request the given route.
	 *
	 * @returns Whether or not the route can be loaded.
	 */
	public async canLoad(): Promise<boolean> {
		return this.auth.fetchCurrentUser();
	}
}
