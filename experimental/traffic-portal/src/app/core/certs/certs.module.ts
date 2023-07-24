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
import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";

import { AppUIModule } from "src/app/app.ui.module";
import { CertViewerComponent } from "src/app/core/certs/cert-viewer/cert-viewer.component";
import { SharedModule } from "src/app/shared/shared.module";

import { CertAuthorComponent } from "./cert-author/cert-author.component";
import { CertDetailComponent } from "./cert-detail/cert-detail.component";

export const ROUTES: Routes = [
	{component: CertViewerComponent, path: "ssl/:xmlId"},
	{component: CertViewerComponent, path: "ssl"}
];

/**
 * Declares the module for SSL certificates. Is seperated since `node-forge` which provides
 * SSL functions is quite large.
 */
@NgModule({
	declarations: [
		CertViewerComponent,
		CertDetailComponent,
		CertAuthorComponent,
	],
	exports: [],
	imports: [
		CommonModule,
		AppUIModule,
		SharedModule,
		RouterModule.forChild(ROUTES)
	]
})
export class CertsModule {
}
