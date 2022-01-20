import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import {
	CacheGroupService,
	CDNService,
	DeliveryServiceService,
	ServerService,
	TypeService
} from "..";
import { CacheGroupService as TestingCacheGroupService } from "./cache-group.service";
import { CDNService as TestingCDNService } from "./cdn.service";
import { DeliveryServiceService as TestingDeliveryServiceService } from "./delivery-service.service";
import { ServerService as TestingServerService } from "./server.service";
import { TypeService as TestingTypeService } from "./type.service";

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
		{provide: ServerService, useClass: TestingServerService},
		{provide: TypeService, useClass: TestingTypeService},
	]
})
export class APITestingModule { }
