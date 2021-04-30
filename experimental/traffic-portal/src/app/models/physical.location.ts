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
 * Represents a Physical Location. Refer to ATC docs for details.
 */
export interface PhysicalLocation {
	address: string;
	city: string;
	comments: string | null;
	email: string | null;
	id?: number;
	lastUpdated?: Date;
	name: string;
	phone: string | null;
	poc: string | null;
	region: string | null;
	regionId: number;
	shortName: string;
	state: string;
	zip: string;
}

export const defaultPhysLoc: PhysicalLocation = {
	address: "",
	city: "",
	comments: null,
	email: null,
	name: "",
	phone: null,
	poc: null,
	region: null,
	regionId: -1,
	shortName: "",
	state: "",
	zip: ""
};
