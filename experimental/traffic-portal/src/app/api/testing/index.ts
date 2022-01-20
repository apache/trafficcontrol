import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { CacheGroupService, CDNService, ServerService } from "..";
import { ServerService as TestingServerService } from "./server.service";
import { CDNService as TestingCDNService } from "./cdn.service";
import { CacheGroupService as TestingCacheGroupService } from "./cache-group.service";

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
		{provide: CacheGroupService, useClass: TestingCacheGroupService},
		{provide: CDNService, useClass: TestingCDNService},
		{provide: ServerService, useClass: TestingServerService}
	]
})
export class APITestingModule { }
