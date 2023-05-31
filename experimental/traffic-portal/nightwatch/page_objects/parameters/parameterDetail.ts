/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { EnhancedPageObject } from "nightwatch";

/**
 * Defines the PageObject for Parameter Details.
 */
export type ParameterDetailPageObject = EnhancedPageObject<{}, typeof parameterDetailPageObject.elements>;

const parameterDetailPageObject = {
	elements: {
		configFile: {
			selector: "input[name='configFile']"
		},
		id: {
			selector: "input[name='id']"
		},
		lastUpdated: {
			selector: "input[name='lastUpdated']"
		},
		name: {
			selector: "input[name='name']"
		},
		saveBtn: {
			selector: "button[type='submit']"
		},
		secure: {
			selector: "input[name='secure']"
		},
		value: {
			selector: "textarea[name='value']"
		}
	},
};

export default parameterDetailPageObject;
