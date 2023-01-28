/*
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
/**
 * This file contains definitons for objects of which Delivery Services are
 * composed - including some convenience functions for the magic numbers used to
 * represent certain properties.
 */

/**
 * DSCapacity represents a response from the API to a request for the capacity
 * of a Delivery Service.
 */
export interface DSCapacity {
	availablePercent: number;
	maintenancePercent: number;
	utilizedPercent: number;
}

/**
 * DSHealth represents a response from the API to a request for the health of a
 * Delivery Service.
 */
export interface DSHealth {
	totalOnline: number;
	totalOffline: number;
}
