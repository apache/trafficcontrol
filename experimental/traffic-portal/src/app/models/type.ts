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
 * This file is for maintaining functionality and modeling related to Apache
 * Traffic Control `Type` objects.
 */

/**
 * Models an arbitrary Type of some object in the database
 */
export interface Type {
	/** A description of the Type. */
	description?: string;
	/** An integral, unique identifier for the Type. */
	id:           number;
	/** The date/time at which the Type was last updated. */
	lastUpdated?: Date;
	/** The Type's name. */
	name:        string;
	/** The database table that uses this Type. */
	useInTable?:  string;
}
