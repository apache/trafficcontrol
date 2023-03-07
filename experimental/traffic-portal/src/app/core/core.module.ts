/**
 * @module src/app/core
 * The "Core" module consists of all TP functionality and components that aren't
 * needed/useful until the user is authenticated.
 *
 * @license Apache-2.0
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
import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { RouterModule, type Routes } from "@angular/router";

import { AppUIModule } from "../app.ui.module";
import { AuthenticatedGuard } from "../guards/authenticated-guard.service";
import { SharedModule } from "../shared/shared.module";

import { CacheGroupDetailsComponent } from "./cache-groups/cache-group-details/cache-group-details.component";
import { CacheGroupTableComponent } from "./cache-groups/cache-group-table/cache-group-table.component";
import { DivisionDetailComponent } from "./cache-groups/divisions/detail/division-detail.component";
import { DivisionsTableComponent } from "./cache-groups/divisions/table/divisions-table.component";
import { RegionDetailComponent } from "./cache-groups/regions/detail/region-detail.component";
import { RegionsTableComponent } from "./cache-groups/regions/table/regions-table.component";
import { ChangeLogsComponent } from "./change-logs/change-logs.component";
import { LastDaysComponent } from "./change-logs/last-days/last-days.component";
import { CurrentuserComponent } from "./currentuser/currentuser.component";
import { UpdatePasswordDialogComponent } from "./currentuser/update-password-dialog/update-password-dialog.component";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { DeliveryserviceComponent } from "./deliveryservice/deliveryservice.component";
import { DsCardComponent } from "./deliveryservice/ds-card/ds-card.component";
import { InvalidationJobsComponent } from "./deliveryservice/invalidation-jobs/invalidation-jobs.component";
import {
	NewInvalidationJobDialogComponent
} from "./deliveryservice/invalidation-jobs/new-invalidation-job-dialog/new-invalidation-job-dialog.component";
import { NewDeliveryServiceComponent } from "./deliveryservice/new-delivery-service/new-delivery-service.component";
import { PhysLocDetailComponent } from "./servers/phys-loc/detail/phys-loc-detail.component";
import { PhysLocTableComponent } from "./servers/phys-loc/table/phys-loc-table.component";
import { ServerDetailsComponent } from "./servers/server-details/server-details.component";
import { ServersTableComponent } from "./servers/servers-table/servers-table.component";
import { UpdateStatusComponent } from "./servers/update-status/update-status.component";
import { StatusDetailsComponent } from "./statuses/status-details/status-details.component";
import { StatusesTableComponent } from "./statuses/statuses-table/statuses-table.component";
import { TypeDetailComponent } from "./types/detail/type-detail.component";
import { TypesTableComponent } from "./types/table/types-table.component";
import { TenantDetailsComponent } from "./users/tenants/tenant-details/tenant-details.component";
import { TenantsComponent } from "./users/tenants/tenants.component";
import { UserDetailsComponent } from "./users/user-details/user-details.component";
import { UserRegistrationDialogComponent } from "./users/user-registration-dialog/user-registration-dialog.component";
import { UsersComponent } from "./users/users.component";

export const ROUTES: Routes = [
	{ component: DashboardComponent, path: "" },
	{ component: DivisionsTableComponent, path: "divisions" },
	{ component: DivisionDetailComponent, path: "divisions/:id" },
	{ component: RegionsTableComponent, path: "regions" },
	{ component: RegionDetailComponent, path: "regions/:id" },
	{ component: UsersComponent, path: "users" },
	{ component: UserDetailsComponent, path: "users/:id"},
	{ component: ServersTableComponent, path: "servers" },
	{ component: ServerDetailsComponent, path: "server/:id" },
	{ component: DeliveryserviceComponent, path: "deliveryservice/:id" },
	{ component: InvalidationJobsComponent, path: "deliveryservice/:id/invalidation-jobs" },
	{ component: CurrentuserComponent, path: "me" },
	{ component: NewDeliveryServiceComponent, path: "new.Delivery.Service" },
	{ component: CacheGroupTableComponent, path: "cache-groups" },
	{ component: CacheGroupDetailsComponent, path: "cache-groups/:id"},
	{ component: TenantsComponent, path: "tenants"},
	{ component: ChangeLogsComponent, path: "change-logs" },
	{ component: TenantDetailsComponent, path: "tenants/:id"},
	{ component: PhysLocDetailComponent, path: "phys-locs/:id" },
	{ component: PhysLocTableComponent, path: "phys-locs" },
	{ component: StatusesTableComponent, path: "statuses" },
	{ component: StatusDetailsComponent, path: "statuses/:id" },
	{ component: TypesTableComponent, path: "types" },
	{ component: TypeDetailComponent, path: "types/:id"},
].map(r => ({...r, canActivate: [AuthenticatedGuard]}));

/**
 * CoreModule contains code that only logged-in users will be served.
 */
@NgModule({
	declarations: [
		UsersComponent,
		ServerDetailsComponent,
		ServersTableComponent,
		DeliveryserviceComponent,
		NewDeliveryServiceComponent,
		CurrentuserComponent,
		UpdatePasswordDialogComponent,
		DashboardComponent,
		DsCardComponent,
		InvalidationJobsComponent,
		CacheGroupTableComponent,
		NewInvalidationJobDialogComponent,
		UpdateStatusComponent,
		UserDetailsComponent,
		TenantsComponent,
		UserRegistrationDialogComponent,
		TenantDetailsComponent,
		ChangeLogsComponent,
		LastDaysComponent,
		UserRegistrationDialogComponent,
		PhysLocTableComponent,
		PhysLocDetailComponent,
		DivisionsTableComponent,
		DivisionDetailComponent,
		RegionsTableComponent,
		RegionDetailComponent,
		CacheGroupDetailsComponent,
		StatusesTableComponent,
		StatusDetailsComponent,
		TypesTableComponent,
		TypeDetailComponent,
	],
	exports: [],
	imports: [
		SharedModule,
		AppUIModule,
		CommonModule,
		RouterModule.forChild(ROUTES)
	]
})
export class CoreModule { }
