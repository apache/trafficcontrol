import { NgModule } from "@angular/core";
import {HTTP_INTERCEPTORS} from "@angular/common/http";
import {CommonModule} from "@angular/common";
import {RouterModule} from "@angular/router";
import {AppUIModule} from "../app.ui.module";
import {DashboardComponent} from "../core/dashboard/dashboard.component";
import {DsCardComponent} from "../core/ds-card/ds-card.component";
import {InvalidationJobsComponent} from "../core/invalidation-jobs/invalidation-jobs.component";
import {CacheGroupTableComponent} from "../core/cache-groups/cache-group-table/cache-group-table.component";
import {NewInvalidationJobDialogComponent} from "../core/invalidation-jobs/new-invalidation-job-dialog/new-invalidation-job-dialog.component";
import {AlertComponent} from "./alert/alert.component";
import {ErrorInterceptor} from "./interceptor/error.interceptor";
import {AlertInterceptor} from "./interceptor/alerts.interceptor";
import {LoadingComponent} from "./loading/loading.component";
import {TpHeaderComponent} from "./tp-header/tp-header.component";
import {GenericTableComponent} from "./generic-table/generic-table.component";
import {LinechartDirective} from "./charts/linechart.directive";
import {AlertService} from "./alert/alert.service";
import {AuthenticationService} from "./authentication/authentication.service";
import {
	CacheGroupService,
	CDNService,
	DeliveryServiceService,
	InvalidationJobService,
	ProfileService,
	ServerService, TypeService, UserService
} from "./api";
import {PhysicalLocationService} from "./api/PhysicalLocationService";
import {CurrentUserService} from "./currentUser/current-user.service";
import {CustomvalidityDirective} from "./validation/customvalidity.directive";
import {OpenableDirective} from "./openable/openable.directive";
import {SSHCellRendererComponent} from "./table-components/ssh-cell-renderer/ssh-cell-renderer.component";
import {UpdateCellRendererComponent} from "./table-components/update-cell-renderer/update-cell-renderer.component";
import {BooleanFilterComponent} from "./table-components/boolean-filter/boolean-filter.component";
import {CurrentuserComponent} from "./currentuser/currentuser.component";
import {UpdatePasswordDialogComponent} from "./currentuser/update-password-dialog/update-password-dialog.component";



/**
 *
 */
@NgModule({
	declarations: [
		AlertComponent,
		LoadingComponent,
		TpHeaderComponent,
		GenericTableComponent,
		BooleanFilterComponent,
		UpdateCellRendererComponent,
		CurrentuserComponent,
		UpdatePasswordDialogComponent,
		DashboardComponent,
		DsCardComponent,
		InvalidationJobsComponent,
		CacheGroupTableComponent,
		NewInvalidationJobDialogComponent,

		CustomvalidityDirective,
		LinechartDirective,
		OpenableDirective
	],
	entryComponents: [
		SSHCellRendererComponent
	],
	exports: [
		AlertComponent,
		LoadingComponent,
		TpHeaderComponent,
		GenericTableComponent,
		BooleanFilterComponent,
		UpdateCellRendererComponent,
		CurrentuserComponent,
		UpdatePasswordDialogComponent,

		CustomvalidityDirective,
		LinechartDirective,
		OpenableDirective
	],
	imports: [
		AppUIModule,
		CommonModule,
		RouterModule
	],
	providers: [
		{multi: true, provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor},
		{multi: true, provide: HTTP_INTERCEPTORS, useClass: AlertInterceptor},
		AlertService,
		AuthenticationService,
		CacheGroupService,
		CDNService,
		CurrentUserService,
		DeliveryServiceService,
		InvalidationJobService,
		PhysicalLocationService,
		ProfileService,
		ServerService,
		TypeService,
		UserService
	],
})
export class SharedModule { }
