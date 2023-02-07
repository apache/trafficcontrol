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

/**
 * Environment is the type of an Angular deployment environment.
 */
export interface Environment {
	/** The version of the Traffic Ops API to be used */
	apiVersion: `${number}.${number}`;
	/** Whether the "Custom" module should be loaded. */
	customModule?: boolean;
	/**
	 * Whether the environment should be treated as a "production" environment.
	 */
	production?: boolean;
	/**
	 * If defined and `true`, the date-reviving HTTP interceptor will attempt to
	 * convert anything that looks like a date into a `Date`. Otherwise, it uses
	 * a specific list of property names known to contain Date information.
	 */
	useExhaustiveDates?: boolean;
}
