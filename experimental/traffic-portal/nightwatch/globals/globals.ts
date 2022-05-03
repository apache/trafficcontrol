/*
*
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
import type {NightwatchBrowser, NightwatchGlobals} from "nightwatch";
import type {CommonPageObject} from "nightwatch/page_objects/common";
import type {DeliveryServiceCardPageObject} from "nightwatch/page_objects/deliveryServiceCard";
import type {DeliveryServiceDetailPageObject} from "nightwatch/page_objects/deliveryServiceDetail";
import type {DeliveryServiceInvalidPageObject} from "nightwatch/page_objects/deliveryServiceInvalidationJobs";
import type {LoginPageObject} from "nightwatch/page_objects/login";
import type {ServersPageObject} from "nightwatch/page_objects/servers";
import type {UsersPageObject} from "nightwatch/page_objects/users";

/**
 * Defines the global nightwatch browser type with our types mixed in.
 */
export type AugmentedBrowser = NightwatchBrowser & {globals: GlobalConfig} & {page:
{
	common: () => CommonPageObject;
	deliveryServiceCard: () => DeliveryServiceCardPageObject;
	deliveryServiceDetail: () => DeliveryServiceDetailPageObject;
	deliveryServiceInvalidationJobs: () => DeliveryServiceInvalidPageObject;
	login: () => LoginPageObject;
	servers: () => ServersPageObject;
	users: () => UsersPageObject;
};};

/**
 * Defines the configuration used for the testing environment
 */
export interface GlobalConfig extends NightwatchGlobals {
	adminPass: string;
	adminUser: string;
	trafficOpsURL: string;
	apiVersion: string;
}
const globals = {
	adminPass: "twelve12",
	adminUser: "admin",
	apiVersion: "4.0",
	trafficOpsURL: "https://localhost:6443"
};

module.exports = globals;
