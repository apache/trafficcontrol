/**
 * @license Apache-2.0
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

import {
	CacheGroupService,
	CDNService,
	ChangeLogsService,
	DeliveryServiceService,
	InvalidationJobService,
	MiscAPIsService,
	OriginService,
	PhysicalLocationService,
	ProfileService,
	ServerService,
	TopologyService,
	TypeService,
	UserService
} from "..";

import { CacheGroupService as TestingCacheGroupService } from "./cache-group.service";
import { CDNService as TestingCDNService } from "./cdn.service";
import { ChangeLogsService as TestingChangeLogsService} from "./change-logs.service";
import { DeliveryServiceService as TestingDeliveryServiceService } from "./delivery-service.service";
import { InvalidationJobService as TestingInvalidationJobService } from "./invalidation-job.service";
import { MiscAPIsService as TestingMiscAPIsService } from "./misc-apis.service";
import { OriginService as TestingOriginService } from "./origin.service";
import { PhysicalLocationService as TestingPhysicalLocationService } from "./physical-location.service";
import { ProfileService as TestingProfileService } from "./profile.service";
import { ServerService as TestingServerService } from "./server.service";
import { TopologyService as TestingTopologyService } from "./topology.service";
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
		{provide: ChangeLogsService, useClass: TestingChangeLogsService},
		{provide: CDNService, useClass: TestingCDNService},
		{provide: DeliveryServiceService, useClass: TestingDeliveryServiceService},
		{provide: InvalidationJobService, useClass: TestingInvalidationJobService},
		{provide: MiscAPIsService, useClass: TestingMiscAPIsService},
		{provide: OriginService, useClass: TestingOriginService},
		{provide: PhysicalLocationService, useClass: TestingPhysicalLocationService},
		{provide: ProfileService, useClass: TestingProfileService},
		{provide: ServerService, useClass: TestingServerService},
		{provide: TopologyService, useClass: TestingTopologyService},
		{provide: TypeService, useClass: TestingTypeService},
		{provide: UserService, useClass: TestingUserService},
		TestingServerService,
		TestingMiscAPIsService
	]
})
export class APITestingModule { }
