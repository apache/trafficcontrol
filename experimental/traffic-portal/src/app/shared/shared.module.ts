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
import { CommonModule, DatePipe } from "@angular/common";
import { HTTP_INTERCEPTORS } from "@angular/common/http";
import { NgModule } from "@angular/core";
import { RouterModule } from "@angular/router";

import { AppUIModule } from "src/app/app.ui.module";
import { DownloadOptionsDialogComponent } from "src/app/shared/generic-table/download-options/download-options-dialog.component";
import { TpHeaderComponent } from "src/app/shared/navigation/tp-header/tp-header.component";
import { TpSidebarComponent } from "src/app/shared/navigation/tp-sidebar/tp-sidebar.component";

import { AlertComponent } from "./alert/alert.component";
import { LinechartDirective } from "./charts/linechart.directive";
import { CollectionChoiceDialogComponent } from "./dialogs/collection-choice-dialog/collection-choice-dialog.component";
import { DecisionDialogComponent } from "./dialogs/decision-dialog/decision-dialog.component";
import { TextDialogComponent } from "./dialogs/text-dialog/text-dialog.component";
import { FileUtilsService } from "./file-utils.service";
import { GenericTableComponent } from "./generic-table/generic-table.component";
import { ImportJsonTxtComponent } from "./import-json-txt/import-json-txt.component";
import { AlertInterceptor } from "./interceptor/alerts.interceptor";
import { DateReviverInterceptor } from "./interceptor/date-reviver.interceptor";
import { ErrorInterceptor } from "./interceptor/error.interceptor";
import { LoadingComponent } from "./loading/loading.component";
import { LoggingService } from "./logging.service";
import { ObscuredTextInputComponent } from "./obscured-text-input/obscured-text-input.component";
import { BooleanFilterComponent } from "./table-components/boolean-filter/boolean-filter.component";
import { EmailCellRendererComponent } from "./table-components/email-cell-renderer/email-cell-renderer.component";
import { SSHCellRendererComponent } from "./table-components/ssh-cell-renderer/ssh-cell-renderer.component";
import { TelephoneCellRendererComponent } from "./table-components/telephone-cell-renderer/telephone-cell-renderer.component";
import { UpdateCellRendererComponent } from "./table-components/update-cell-renderer/update-cell-renderer.component";
import { TreeSelectComponent } from "./tree-select/tree-select.component";
import { CustomvalidityDirective } from "./validation/customvalidity.directive";

/**
 * SharedModule contains common code that modules can import independently.
 */
@NgModule({
	declarations: [
		AlertComponent,
		LoadingComponent,
		TpHeaderComponent,
		TpSidebarComponent,
		GenericTableComponent,
		BooleanFilterComponent,
		UpdateCellRendererComponent,
		CustomvalidityDirective,
		LinechartDirective,
		SSHCellRendererComponent,
		EmailCellRendererComponent,
		TelephoneCellRendererComponent,
		ObscuredTextInputComponent,
		TreeSelectComponent,
		TextDialogComponent,
		DecisionDialogComponent,
		CollectionChoiceDialogComponent,
		ImportJsonTxtComponent,
		DownloadOptionsDialogComponent
	],
	exports: [
		AlertComponent,
		LoadingComponent,
		TpHeaderComponent,
		TpSidebarComponent,
		GenericTableComponent,
		BooleanFilterComponent,
		UpdateCellRendererComponent,
		CustomvalidityDirective,
		LinechartDirective,
		ObscuredTextInputComponent,
		TreeSelectComponent
	],
	imports: [
		AppUIModule,
		CommonModule,
		RouterModule
	],
	providers: [
		{ multi: true, provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor },
		{ multi: true, provide: HTTP_INTERCEPTORS, useClass: AlertInterceptor },
		{ multi: true, provide: HTTP_INTERCEPTORS, useClass: DateReviverInterceptor },
		FileUtilsService,
		DatePipe,
		LoggingService
	]
})
export class SharedModule { }
