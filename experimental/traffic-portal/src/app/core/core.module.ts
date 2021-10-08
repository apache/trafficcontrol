import { NgModule } from "@angular/core";
import {RouterModule, Routes} from "@angular/router";
import {CommonModule} from "@angular/common";
import {AppUIModule} from "../app.ui.module";
import {SharedModule} from "../shared/shared.module";
import {AuthenticatedGuard} from "../guards/authenticated-guard.service";
import {InvalidationJobsComponent} from "./invalidation-jobs/invalidation-jobs.component";
import {UsersComponent} from "./users/users.component";
import {ServerDetailsComponent} from "./servers/server-details/server-details.component";
import {ServersTableComponent} from "./servers/servers-table/servers-table.component";
import {UpdateStatusComponent} from "./servers/update-status/update-status.component";
import {DeliveryserviceComponent} from "./deliveryservice/deliveryservice.component";
import {NewDeliveryServiceComponent} from "./new-delivery-service/new-delivery-service.component";
import {DashboardComponent} from "./dashboard/dashboard.component";
import {CacheGroupTableComponent} from "./cache-groups/cache-group-table/cache-group-table.component";
import {CurrentuserComponent} from "./currentuser/currentuser.component";
import {UpdatePasswordDialogComponent} from "./currentuser/update-password-dialog/update-password-dialog.component";
import {DsCardComponent} from "./ds-card/ds-card.component";
import {NewInvalidationJobDialogComponent} from "./invalidation-jobs/new-invalidation-job-dialog/new-invalidation-job-dialog.component";


const routes: Routes = [
	{ component: DashboardComponent, path: "", canActivate: [AuthenticatedGuard]},
	{ component: UsersComponent, path: "users", canActivate: [AuthenticatedGuard]},
	{ component: ServersTableComponent, path: "servers", canActivate: [AuthenticatedGuard]},
	{ component: ServerDetailsComponent, path: "server/:id", canActivate: [AuthenticatedGuard] },
	{ component: DeliveryserviceComponent, path: "deliveryservice/:id", canActivate: [AuthenticatedGuard] },
	{ component: InvalidationJobsComponent, path: "deliveryservice/:id/invalidation-jobs", canActivate: [AuthenticatedGuard] },
	{ component: CurrentuserComponent, path: "me", canActivate: [AuthenticatedGuard] },
	{ component: NewDeliveryServiceComponent, path: "new.Delivery.Service", canActivate: [AuthenticatedGuard] },
	{ component: CacheGroupTableComponent, path: "cache-groups", canActivate: [AuthenticatedGuard] }
];

/**
 *
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
		UpdateStatusComponent
	],
	exports: [
	],
	imports: [
		SharedModule,
		AppUIModule,
		CommonModule,
		RouterModule.forChild(routes)
	]
})
export class CoreModule { }
