/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

*/

/**
 * PropertiesModel models the "properties" (i.e. configuration) of Traffic
 * Portal.
 */
class PropertiesModel {

	/**
	 * @type {Partial<typeof import("../../traffic_portal_properties.json").properties>}
	 */
	properties = {};
	loaded = false;

	/**
	 * Loads in properties if not already done.
	 *
	 * @param {Partial<typeof import("../../traffic_portal_properties.json").properties>} properties
	 */
	setProperties(properties) {
		if (this.loaded) return;
		this.properties = properties;
		this.loaded = true;
	};

};

PropertiesModel.$inject = [];
module.exports = PropertiesModel;
