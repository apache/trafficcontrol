import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { CDNService } from "./cdn.service";
import { ServerService } from "./server.service";
import { CacheGroupService } from "./cache-group.service";

export * from "./cache-group.service";
export * from "./cdn.service";
export * from "./server.service";

/**
 * The API Module contains all logic used to access the Traffic Ops API.
 */
@NgModule({
	declarations: [],
	imports: [
		CommonModule
	],
	providers: [
		CacheGroupService,
		CDNService,
		ServerService
	]
})
export class APIModule { }
