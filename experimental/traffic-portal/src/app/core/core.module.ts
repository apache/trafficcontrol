import { NgModule } from "@angular/core";
import {RouterModule, Routes} from "@angular/router";
import {CommonModule} from "@angular/common";
import {AppUIModule} from "../app.ui.module";
import {SharedModule} from "../shared/shared.module";
import {CurrentuserComponent} from "../shared/currentuser/currentuser.component";
import {AuthenticationGuard} from "../authentication.guard";
import {InvalidationJobsComponent} from "./invalidation-jobs/invalidation-jobs.component";
import {UsersComponent} from "./users/users.component";
import {ServerDetailsComponent} from "./servers/server-details/server-details.component";
import {ServersTableComponent} from "./servers/servers-table/servers-table.component";
import {UpdateStatusComponent} from "./servers/update-status/update-status.component";
import {DeliveryserviceComponent} from "./deliveryservice/deliveryservice.component";
import {NewDeliveryServiceComponent} from "./new-delivery-service/new-delivery-service.component";
import {DashboardComponent} from "./dashboard/dashboard.component";
import {CacheGroupTableComponent} from "./cache-groups/cache-group-table/cache-group-table.component";


const routes: Routes = [
	{ component: DashboardComponent, path: "", canActivate: [AuthenticationGuard]},
	{ component: UsersComponent, path: "users", canActivate: [AuthenticationGuard], data:{animation:"users"}},
	{ component: ServersTableComponent, path: "servers", canActivate: [AuthenticationGuard], data:{animation:"servers"}},
	{ component: ServerDetailsComponent, path: "server/:id", canActivate: [AuthenticationGuard] },
	{ component: DeliveryserviceComponent, path: "deliveryservice/:id", canActivate: [AuthenticationGuard] },
	{ component: InvalidationJobsComponent, path: "deliveryservice/:id/invalidation-jobs", canActivate: [AuthenticationGuard] },
	{ component: CurrentuserComponent, path: "me", canActivate: [AuthenticationGuard] },
	{ component: NewDeliveryServiceComponent, path: "new.Delivery.Service", canActivate: [AuthenticationGuard] },
	{ component: CacheGroupTableComponent, path: "cache-groups", canActivate: [AuthenticationGuard] }
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
		UpdateStatusComponent
	],
	exports: [
		UsersComponent,
		ServerDetailsComponent,
		ServersTableComponent,
		DeliveryserviceComponent,
		NewDeliveryServiceComponent,
		UpdateStatusComponent
	],
	imports: [
		SharedModule,
		AppUIModule,
		CommonModule,
		RouterModule.forChild(routes)
	]
})
export class CoreModule { }
