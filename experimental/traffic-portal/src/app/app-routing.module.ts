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
import {NgModule} from "@angular/core";
import {RouterModule, Routes} from "@angular/router";

import {LoginComponent} from "./login/login.component";
import {DemoComponent} from "./demo/demo.component";

const routes: Routes = [
	{component: LoginComponent, path: "login"},
	{
		children: [{
			loadChildren: async () => import("./core/core.module")
				.then(mod => mod.CoreModule),
			path: ""
		}],
		path: "core"
	},
	{ component: DemoComponent, path: "demo" }
];

/**
 * AppRoutingModule provides routing configuration for the app.
 */
@NgModule({
	exports: [RouterModule],
	imports: [RouterModule.forRoot(routes, {
		initialNavigation: "enabled",
		relativeLinkResolution: "legacy"
	})],
})
// This is a necessary empty class. All of its data/logic come from the decorator.
// eslint-disable-next-line @typescript-eslint/no-extraneous-class
export class AppRoutingModule { }
