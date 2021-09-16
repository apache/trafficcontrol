import { NgModule } from "@angular/core";
import {RouterModule, Routes} from "@angular/router";
import {CommonModule} from "@angular/common";
import {AppUIModule} from "../app.ui.module";
import {SharedModule} from "../shared/shared.module";
import {CurrentuserComponent} from "../shared/currentuser/currentuser.component";
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
	{ component: DashboardComponent, path: ""},
	{ component: UsersComponent, path: "users" },
	{ component: ServersTableComponent, path: "servers" },
	{ component: ServerDetailsComponent, path: "server/:id" },
	{ component: DeliveryserviceComponent, path: "deliveryservice/:id" },
	{ component: InvalidationJobsComponent, path: "deliveryservice/:id/invalidation-jobs" },
	{component: CurrentuserComponent, path: "me"},
	{component: NewDeliveryServiceComponent, path: "new.Delivery.Service"},
	{component: CacheGroupTableComponent, path: "cache-groups"}
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
