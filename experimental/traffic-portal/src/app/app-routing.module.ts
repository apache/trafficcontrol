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
import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";

import { CurrentuserComponent } from "./components/currentuser/currentuser.component";
import { DashboardComponent } from "./components/dashboard/dashboard.component";
import { DeliveryserviceComponent } from "./components/deliveryservice/deliveryservice.component";
import { InvalidationJobsComponent } from "./components/invalidation-jobs/invalidation-jobs.component";
import { LoginComponent } from "./components/login/login.component";
import { NewDeliveryServiceComponent } from "./components/new-delivery-service/new-delivery-service.component";
import { ServersTableComponent } from "./components/servers/servers-table/servers-table.component";
import { UsersComponent } from "./components/users/users.component";

const routes: Routes = [
	{ path: "", component: DashboardComponent },
	{ path: "login", component: LoginComponent },
	{ path: "users", component: UsersComponent},
	{ path: "me", component: CurrentuserComponent},
	{ path: "new.Delivery.Service", component: NewDeliveryServiceComponent},
	{ path: "deliveryservice/:id", component: DeliveryserviceComponent},
	{ path: "deliveryservice/:id/invalidation-jobs", component: InvalidationJobsComponent},
	{ path: "servers", component: ServersTableComponent},
];

/**
 * AppRoutingModule provides routing configuration for the app.
 */
@NgModule({
	exports: [RouterModule],
	imports: [RouterModule.forRoot(routes, {
		initialNavigation: "enabled"
})],
})
export class AppRoutingModule { }
