import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { CDNService, ServerService } from "..";
import { ServerService as TestingServerService } from "./server.service";
import { CDNService as TestingCDNService } from "./cdn.service";

export * from "./cdn.service";
export * from "./server.service";

/**
 * The API Testing Module provides mock services that allow components to use
 * the Traffic Ops API without actually requiring a running Traffic Ops.
 */
@NgModule({
	declarations: [],
	imports: [
		CommonModule
	],
	providers: [
		{provide: CDNService, useClass: TestingCDNService},
		{provide: ServerService, useClass: TestingServerService}
	]
})
export class APITestingModule { }
