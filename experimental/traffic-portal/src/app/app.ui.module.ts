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
import { AgGridModule } from "ag-grid-angular";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {FontAwesomeModule} from "@fortawesome/angular-fontawesome";
import {MatButtonModule} from "@angular/material/button";
import {MatCardModule} from "@angular/material/card";
import {MatDividerModule} from "@angular/material/divider";
import {MatInputModule} from "@angular/material/input";
import {MatListModule} from "@angular/material/list";
import {MatRadioModule} from "@angular/material/radio";
import {MatSnackBarModule} from "@angular/material/snack-bar";
import {MatStepperModule} from "@angular/material/stepper";
import {MatNativeDateModule} from "@angular/material/core";
import {MatDialogModule} from "@angular/material/dialog";
import {MatDatepickerModule} from "@angular/material/datepicker";
import {MatToolbarModule} from "@angular/material/toolbar";
import {MatExpansionModule} from "@angular/material/expansion";
import {MatButtonToggleModule} from "@angular/material/button-toggle";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

/**
 * AppUIModule is the Angular Module that contains the ui dependencies of
 * the app.
 */
@NgModule({
	bootstrap: [],
	exports: [
		AgGridModule,

		BrowserAnimationsModule,

		ReactiveFormsModule,
		FormsModule,

		FontAwesomeModule,

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
		MatButtonToggleModule
	],
})
export class AppUIModule {}
