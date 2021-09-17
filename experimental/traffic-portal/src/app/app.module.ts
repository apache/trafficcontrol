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

/**
 * This file contains the definition for the entire app. Its syntax is a bit arcane, but hopefully
 * by copy/pasting any novice can add a new component - though honestly you should just use
 * `ng generate` to create new things (and then fix formatting/missing license)
 */

import { HttpClientModule, HTTP_INTERCEPTORS } from "@angular/common/http";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule, FormsModule } from "@angular/forms";
import { BrowserModule } from "@angular/platform-browser";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatNativeDateModule } from "@angular/material/core";
import { MatDividerModule } from "@angular/material/divider";
import { MatExpansionModule } from "@angular/material/expansion";
import { MatInputModule } from "@angular/material/input";
import { MatListModule } from "@angular/material/list";
import { MatRadioModule } from "@angular/material/radio";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatStepperModule } from "@angular/material/stepper";
import { MatToolbarModule } from "@angular/material/toolbar";
import { MatDialogModule } from "@angular/material/dialog";
import { MatDatepickerModule } from "@angular/material/datepicker";

import { FontAwesomeModule } from "@fortawesome/angular-fontawesome";
import { AgGridModule } from "ag-grid-angular";
import * as Chart from "chart.js";

// Routing, Components, Directives and Interceptors
import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { AlertComponent } from "./components/alert/alert.component";
import { CurrentuserComponent } from "./components/currentuser/currentuser.component";
import { DashboardComponent } from "./components/dashboard/dashboard.component";
import { DeliveryserviceComponent } from "./components/deliveryservice/deliveryservice.component";
import { DsCardComponent } from "./components/ds-card/ds-card.component";
import { InvalidationJobsComponent } from "./components/invalidation-jobs/invalidation-jobs.component";
import { LoadingComponent } from "./components/loading/loading.component";
import { LoginComponent } from "./components/login/login.component";
import { NewDeliveryServiceComponent } from "./components/new-delivery-service/new-delivery-service.component";
import { ServersTableComponent } from "./components/servers/servers-table/servers-table.component";
import { SSHCellRendererComponent } from "./components/table-components/ssh-cell-renderer/ssh-cell-renderer.component";
import { TpHeaderComponent } from "./components/tp-header/tp-header.component";
import { UsersComponent } from "./components/users/users.component";
import { CustomvalidityDirective } from "./directives/customvalidity.directive";
import { LinechartDirective } from "./directives/linechart.directive";
import { OpenableDirective } from "./directives/openable.directive";
import { AlertInterceptor } from "./interceptor/alerts.interceptor";
import { ErrorInterceptor } from "./interceptor/error.interceptor";
import { GenericTableComponent } from "./components/generic-table/generic-table.component";
import { CacheGroupTableComponent } from "./components/cache-groups/cache-group-table/cache-group-table.component";
import { BooleanFilterComponent } from "./components/table-components/boolean-filter/boolean-filter.component";
import { ServerDetailsComponent } from "./components/servers/server-details/server-details.component";
import { UpdateCellRendererComponent } from "./components/table-components/update-cell-renderer/update-cell-renderer.component";
import { UpdateStatusComponent } from "./components/servers/update-status/update-status.component";
import {
	NewInvalidationJobDialogComponent
} from "./components/invalidation-jobs/new-invalidation-job-dialog/new-invalidation-job-dialog.component";
import { UpdatePasswordDialogComponent } from "./components/currentuser/update-password-dialog/update-password-dialog.component";
import { ChartsComponent } from './components/charts/charts.component';

// TODO: Figure out the actual typing here.
Chart.plugins.register({
	beforeDraw: (chartInstance: Chart & {
		chart: {
			ctx: {
				fillStyle: string;
				fillRect: (a: number, b: number, c: number, d: number) => void;
			};
			width: number;
			height: number;
		};
	}) => {
		const ctx = chartInstance.chart.ctx;
		ctx.fillStyle = "white";
		ctx.fillRect(0, 0, chartInstance.chart.width, chartInstance.chart.height);
	},
	id: "whiteBackground",
});

/**
 * AppModule is the single Angular Module that contains the entire
 * front-end part of the app (which is nearly all of it).
 */
@NgModule({
	bootstrap: [AppComponent],
	declarations: [
		AppComponent,
		LoginComponent,
		DashboardComponent,
		DsCardComponent,
		AlertComponent,
		UsersComponent,
		NewDeliveryServiceComponent,
		TpHeaderComponent,
		LoadingComponent,
		DeliveryserviceComponent,
		LinechartDirective,
		InvalidationJobsComponent,
		OpenableDirective,
		CustomvalidityDirective,
		CurrentuserComponent,
		ServersTableComponent,
		GenericTableComponent,
		CacheGroupTableComponent,
		BooleanFilterComponent,
		ServerDetailsComponent,
		UpdateCellRendererComponent,
		UpdateStatusComponent,
		NewInvalidationJobDialogComponent,
		UpdatePasswordDialogComponent,
  		ChartsComponent
	],
	entryComponents: [
		SSHCellRendererComponent
	],
	imports: [
		BrowserModule.withServerTransition({ appId: "serverApp" }),
		AppRoutingModule,
		HttpClientModule,
		ReactiveFormsModule,
		FormsModule,
		FontAwesomeModule,
		AgGridModule.withComponents([]),
		BrowserAnimationsModule,
		MatButtonModule,
		MatCardModule,
		MatDividerModule,
		MatExpansionModule,
		MatInputModule,
		MatListModule,
		MatRadioModule,
		MatSnackBarModule,
		MatStepperModule,
		MatToolbarModule,
		MatDialogModule,
		MatDatepickerModule,
		MatNativeDateModule,
	],
	providers: [
		{multi: true, provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor},
		{multi: true, provide: HTTP_INTERCEPTORS, useClass: AlertInterceptor}
	],
})
// This is a necessary empty class. All of its data/logic come from the decorator.
// eslint-disable-next-line @typescript-eslint/no-extraneous-class
export class AppModule { }
