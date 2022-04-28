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
import {type NightwatchGlobals} from "nightwatch";

/**
 * Defines the configuration used for the testing environment
 */
export interface GlobalConfig extends NightwatchGlobals {
	adminPass: string;
	adminUser: string;
	trafficOpsURL: string;
}
const config = {
	adminPass: "twelve12",
	adminUser: "admin",
	trafficOpsURL: "https://localhost:6443"
};

module.exports = config;
