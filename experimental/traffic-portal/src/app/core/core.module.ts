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

import { DivisionDetailComponent } from "src/app/core/cache-groups/divisions/detail/division-detail.component";
import { DivisionsTableComponent } from "src/app/core/cache-groups/divisions/table/divisions-table.component";
import { RegionDetailComponent } from "src/app/core/cache-groups/regions/detail/region-detail.component";
import { RegionsTableComponent } from "src/app/core/cache-groups/regions/table/regions-table.component";
import { LastDaysComponent } from "src/app/core/change-logs/last-days/last-days.component";

import { AppUIModule } from "../app.ui.module";
import { AuthenticatedGuard } from "../guards/authenticated-guard.service";
import { SharedModule } from "../shared/shared.module";

import { CacheGroupTableComponent } from "./cache-groups/cache-group-table/cache-group-table.component";
import { DiffFieldComponent } from "./cdns/snapshot/diff-field/diff-field.component";
import { ContentRouterComponent } from "./cdns/snapshot/router-diff/content-router/content-router.component";
import { RouterDiffComponent } from "./cdns/snapshot/router-diff/router-diff.component";
import { ContentServerComponent } from "./cdns/snapshot/server-diff/content-server/content-server.component";
import { ServerDiffComponent } from "./cdns/snapshot/server-diff/server-diff.component";
import { ChangeLogsComponent } from "./change-logs/change-logs.component";
import { CurrentuserComponent } from "./currentuser/currentuser.component";
import { UpdatePasswordDialogComponent } from "./currentuser/update-password-dialog/update-password-dialog.component";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { DeliveryserviceComponent } from "./deliveryservice/deliveryservice.component";
import { DsCardComponent } from "./ds-card/ds-card.component";
import { InvalidationJobsComponent } from "./invalidation-jobs/invalidation-jobs.component";
import {
	NewInvalidationJobDialogComponent
} from "./invalidation-jobs/new-invalidation-job-dialog/new-invalidation-job-dialog.component";
import { NewDeliveryServiceComponent } from "./new-delivery-service/new-delivery-service.component";
import { ServerDetailsComponent } from "./servers/server-details/server-details.component";
import { ServersTableComponent } from "./servers/servers-table/servers-table.component";
import { UpdateStatusComponent } from "./servers/update-status/update-status.component";
import { TenantDetailsComponent } from "./users/tenants/tenant-details/tenant-details.component";
import { TenantsComponent } from "./users/tenants/tenants.component";
import { UserDetailsComponent } from "./users/user-details/user-details.component";
import { UserRegistrationDialogComponent } from "./users/user-registration-dialog/user-registration-dialog.component";
import { UsersComponent } from "./users/users.component";

export const ROUTES: Routes = [
	{ canActivate: [AuthenticatedGuard], component: DashboardComponent, path: "" },
	{ canActivate: [AuthenticatedGuard], component: DivisionsTableComponent, path: "divisions" },
	{ canActivate: [AuthenticatedGuard], component: DivisionDetailComponent, path: "division/:id" },
	{ canActivate: [AuthenticatedGuard], component: RegionsTableComponent, path: "regions" },
	{ canActivate: [AuthenticatedGuard], component: RegionDetailComponent, path: "region/:id" },
	{ canActivate: [AuthenticatedGuard], component: UsersComponent, path: "users" },
	{ canActivate: [AuthenticatedGuard], component: UserDetailsComponent, path: "users/:id"},
	{ canActivate: [AuthenticatedGuard], component: ServersTableComponent, path: "servers" },
	{ canActivate: [AuthenticatedGuard], component: ServerDetailsComponent, path: "server/:id" },
	{ canActivate: [AuthenticatedGuard], component: DeliveryserviceComponent, path: "deliveryservice/:id" },
	{ canActivate: [AuthenticatedGuard], component: InvalidationJobsComponent, path: "deliveryservice/:id/invalidation-jobs" },
	{ canActivate: [AuthenticatedGuard], component: CurrentuserComponent, path: "me" },
	{ canActivate: [AuthenticatedGuard], component: NewDeliveryServiceComponent, path: "new.Delivery.Service" },
	{ canActivate: [AuthenticatedGuard], component: CacheGroupTableComponent, path: "cache-groups" },
	{ canActivate: [AuthenticatedGuard], component: TenantsComponent, path: "tenants"},
	{ canActivate: [AuthenticatedGuard], component: ChangeLogsComponent, path: "change-logs" },
	{ canActivate: [AuthenticatedGuard], component: TenantDetailsComponent, path: "tenants/:id"}
];

/**
 * CoreModule contains code that only logged in users will be served.
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
		DivisionsTableComponent,
		DivisionDetailComponent,
		RegionsTableComponent,
		RegionDetailComponent,
		DiffFieldComponent,
		ContentRouterComponent,
		ContentServerComponent,
		RouterDiffComponent,
		ServerDiffComponent,
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
