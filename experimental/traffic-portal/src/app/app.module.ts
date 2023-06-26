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

import { HttpClientModule } from "@angular/common/http";
import { NgModule } from "@angular/core";
import { BrowserModule, provideClientHydration } from "@angular/platform-browser";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import * as Chart from "chart.js";

// Routing, Components, Directives and Interceptors
import { APIModule } from "./api";
import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { AppUIModule } from "./app.ui.module";
import { AuthenticatedGuard } from "./guards/authenticated-guard.service";
import { LoginComponent } from "./login/login.component";
import { ResetPasswordDialogComponent } from "./login/reset-password-dialog/reset-password-dialog.component";
import { SharedModule } from "./shared/shared.module";

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
		ResetPasswordDialogComponent
	],
	imports: [
		BrowserModule.withServerTransition({appId: "serverApp"}),
		BrowserAnimationsModule,
		AppRoutingModule,
		HttpClientModule,
		AppUIModule,
		SharedModule,
		APIModule
	],
	providers: [
		AuthenticatedGuard,
		provideClientHydration()
	]
})
export class AppModule { }
