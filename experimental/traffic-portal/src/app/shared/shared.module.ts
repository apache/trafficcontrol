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
import { HTTP_INTERCEPTORS } from "@angular/common/http";
import { CommonModule } from "@angular/common";
import { RouterModule } from "@angular/router";

import { AppUIModule } from "src/app/app.ui.module";

import { AlertComponent } from "./alert/alert.component";
import { AlertInterceptor } from "./interceptor/alerts.interceptor";
import { AlertService } from "./alert/alert.service";
import { BooleanFilterComponent } from "./table-components/boolean-filter/boolean-filter.component";
import { CurrentUserService } from "./currentUser/current-user.service";
import { CustomvalidityDirective } from "./validation/customvalidity.directive";
import { ErrorInterceptor } from "./interceptor/error.interceptor";
import { GenericTableComponent } from "./generic-table/generic-table.component";
import { LinechartDirective } from "./charts/linechart.directive";
import { LoadingComponent } from "./loading/loading.component";
import { SSHCellRendererComponent } from "./table-components/ssh-cell-renderer/ssh-cell-renderer.component";
import { TpHeaderComponent } from "./tp-header/tp-header.component";
import { UpdateCellRendererComponent } from "./table-components/update-cell-renderer/update-cell-renderer.component";

/**
 * SharedModule contains common code that modules can import independently.
 */
@NgModule({
	declarations: [
		AlertComponent,
		LoadingComponent,
		TpHeaderComponent,
		GenericTableComponent,
		BooleanFilterComponent,
		UpdateCellRendererComponent,
		CustomvalidityDirective,
		LinechartDirective,
		SSHCellRendererComponent,
	],
	exports: [
		AlertComponent,
		LoadingComponent,
		TpHeaderComponent,
		GenericTableComponent,
		BooleanFilterComponent,
		UpdateCellRendererComponent,
		CustomvalidityDirective,
		LinechartDirective,
	],
	imports: [
		AppUIModule,
		CommonModule,
		RouterModule
	],
	providers: [
		{ multi: true, provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor },
		{ multi: true, provide: HTTP_INTERCEPTORS, useClass: AlertInterceptor },
		AlertService,
		CurrentUserService,
	]
})
export class SharedModule { }
