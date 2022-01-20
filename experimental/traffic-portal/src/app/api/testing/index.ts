import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import {
	CacheGroupService,
	CDNService,
	DeliveryServiceService,
	InvalidationJobService,
	ProfileService,
	ServerService,
	TypeService,
	UserService
} from "..";
import { CacheGroupService as TestingCacheGroupService } from "./cache-group.service";
import { CDNService as TestingCDNService } from "./cdn.service";
import { DeliveryServiceService as TestingDeliveryServiceService } from "./delivery-service.service";
import { InvalidationJobService as TestingInvalidationJobService } from "./invalidation-job.service";
import { ProfileService as TestingProfileService } from "./profile.service";
import { ServerService as TestingServerService } from "./server.service";
import { TypeService as TestingTypeService } from "./type.service";
import { UserService as TestingUserService } from "./user.service";

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
		{provide: DeliveryServiceService, useClass: TestingDeliveryServiceService},
		{provide: InvalidationJobService, useClass: TestingInvalidationJobService},
		{provide: ProfileService, useClass: TestingProfileService},
		{provide: ServerService, useClass: TestingServerService},
		{provide: TypeService, useClass: TestingTypeService},
		{provide: UserService, useClass: TestingUserService}
	]
})
export class APITestingModule { }
