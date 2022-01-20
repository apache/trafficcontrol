import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { CDNService } from "./cdn.service";
import { CacheGroupService } from "./cache-group.service";
import { DeliveryServiceService } from "./delivery-service.service";
import { InvalidationJobService } from "./invalidation-job.service";
import { PhysicalLocationService } from "./physical-location.service";
import { ProfileService } from "./profile.service";
import { ServerService } from "./server.service";
import { TypeService } from "./type.service";
import { UserService } from "./user.service";

export * from "./cache-group.service";
export * from "./cdn.service";
export * from "./delivery-service.service";
export * from "./invalidation-job.service";
export * from "./physical-location.service";
export * from "./profile.service";
export * from "./server.service";
export * from "./type.service";
export * from "./user.service";

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
		DeliveryServiceService,
		InvalidationJobService,
		PhysicalLocationService,
		ProfileService,
		ServerService,
		TypeService,
		UserService,
	]
})
export class APIModule { }
