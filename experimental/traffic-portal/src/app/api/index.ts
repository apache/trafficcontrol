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

import { ChangeLogsService } from "src/app/api/change-logs.service";

import { CacheGroupService } from "./cache-group.service";
import { CDNService } from "./cdn.service";
import { DeliveryServiceService } from "./delivery-service.service";
import { InvalidationJobService } from "./invalidation-job.service";
import { MiscAPIsService } from "./misc-apis.service";
import { OriginService } from "./origin.service";
import { PhysicalLocationService } from "./physical-location.service";
import { ProfileService } from "./profile.service";
import { ServerService } from "./server.service";
import { TopologyService } from "./topology.service";
import { TypeService } from "./type.service";
import { UserService } from "./user.service";

export * from "./cache-group.service";
export * from "./cdn.service";
export * from "./change-logs.service";
export * from "./delivery-service.service";
export * from "./invalidation-job.service";
export * from "./misc-apis.service";
export * from "./physical-location.service";
export * from "./profile.service";
export * from "./server.service";
export * from "./topology.service";
export * from "./type.service";
export * from "./user.service";
export * from "./origin.service";

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
		ChangeLogsService,
		DeliveryServiceService,
		InvalidationJobService,
		MiscAPIsService,
		PhysicalLocationService,
		ProfileService,
		ServerService,
		TopologyService,
		TypeService,
		UserService,
		OriginService,
	]
})
export class APIModule { }
