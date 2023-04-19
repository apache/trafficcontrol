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
import { NgModule, type Type } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";

import { type CoreModule } from "src/app/core/core.module";
import { type CustomModule } from "src/app/custom/custom.module";
import { environment } from "src/environments/environment";

import { AuthenticatedGuard } from "./guards/authenticated-guard.service";
import { LoginComponent } from "./login/login.component";

const routes: Routes = [
	{component: LoginComponent, path: "login"},
	{
		canLoad: [AuthenticatedGuard],
		children: [{
			loadChildren: async (): Promise<Type<CoreModule>> => import("./core/core.module")
				.then(mod => mod.CoreModule),
			path: ""
		}],
		path: "core",
	},
	{path: "", pathMatch: "full", redirectTo: "login"}
];

if (environment.customModule) {
	routes.push({
		children: [{
			loadChildren: async (): Promise<Type<CustomModule>> =>
				import("./custom/custom.module").then(mod => mod.CustomModule),
			path: ""
		}],
		path: "custom"
	});
}

/**
 * AppRoutingModule provides routing configuration for the app.
 */
@NgModule({
	exports: [RouterModule],
	imports: [RouterModule.forRoot(routes, {
		initialNavigation: "enabledNonBlocking"
	})],
})
export class AppRoutingModule { }
