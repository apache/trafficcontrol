import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";

import { AppUIModule } from "src/app/app.ui.module";
import { CertViewerComponent } from "src/app/core/certs/cert-viewer/cert-viewer.component";
import { SharedModule } from "src/app/shared/shared.module";

import { CertAuthorComponent } from "./cert-author/cert-author.component";
import { CertDetailComponent } from "./cert-detail/cert-detail.component";

export const ROUTES: Routes = [
	{component: CertViewerComponent, path: "ssl/ds/:xmlId"},
	{component: CertViewerComponent, path: "ssl"}
];

/**
 *
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
